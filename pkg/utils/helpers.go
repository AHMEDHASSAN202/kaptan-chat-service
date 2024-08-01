package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime/debug"
	"samm/pkg/logger"
	"strings"
	"time"
)

const DefaultTimeFormat = "15:04:05"
const DefaultHourTimeFormat = "15:04"
const DefaultDateFormat = "2006-01-02"

// RemoveDuplicates removes duplicate values from a slice.
// T must be a comparable type.
func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	result := []T{}
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func ConvertStringIdsToObjectIds(ids []string) []primitive.ObjectID {
	var _ids []primitive.ObjectID
	for _, id := range ids {
		_ids = append(_ids, ConvertStringIdToObjectId(id))
	}
	return _ids
}

func ConvertObjectIdsToStringIds(ids []primitive.ObjectID) []string {
	var _ids []string
	for _, id := range ids {
		_ids = append(_ids, ConvertObjectIdToStringId(id))
	}
	return _ids
}

func ConvertStringIdToObjectId(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return objectId
}

func IsValidateObjectId(id string) bool {
	return ConvertStringIdToObjectId(id) != primitive.NilObjectID
}

func IsObjectIdValid(id primitive.ObjectID) bool {
	return id != primitive.NilObjectID && id.Hex() != "" && primitive.IsValidObjectID(id.Hex())
}

func ConvertObjectIdToStringId(id primitive.ObjectID) string {
	return id.Hex()
}

// Contains checks if a value exists in a slice.
func Contains(slice interface{}, value interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	// Check if the provided slice is actually a slice
	if sliceValue.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

// ContainsAny checks if any value in slice1 is present in slice2.
func ContainsAny(slice1 interface{}, slice2 interface{}) bool {
	slice1Value := reflect.ValueOf(slice1)
	// Check if the provided slice1 is actually a slice
	if slice1Value.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < slice1Value.Len(); i++ {
		if Contains(slice2, slice1Value.Index(i).Interface()) {
			return true
		}
	}
	return false
}

func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func ConvertArrStructToInterfaceArr[T any](t []T) []interface{} {
	s := make([]interface{}, len(t))
	for i, v := range t {
		s[i] = v
	}
	return s
}

func ValidateIDsIsMongoObjectIds(fl validator.FieldLevel) bool {
	entityIDs := fl.Field().Interface().([]string)
	if len(entityIDs) == 0 {
		return true
	}
	for _, id := range entityIDs {
		if !IsValidateObjectId(id) {
			return false
		}
	}
	return true
}

func ValidateCountryIds(fl validator.FieldLevel) bool {
	countryIDs := fl.Field().Interface().([]string)
	if len(countryIDs) == 0 {
		return false
	}
	for _, country := range countryIDs {
		if !Contains(Countries, strings.ToUpper(country)) {
			return false
		}
	}
	return true
}

// DiffStructs returns a map of field names and their differing values between two structs.
func DiffStructs(a, b interface{}) []string {
	differences := make([]string, 0)
	compareStructs(reflect.ValueOf(a), reflect.ValueOf(b), "", differences)
	return differences
}

// DiffStructs returns differing keys values between two structs.
func compareStructs(valA, valB reflect.Value, parentField string, differences []string) {
	if valA.Kind() == reflect.Ptr {
		valA = valA.Elem()
	}
	if valB.Kind() == reflect.Ptr {
		valB = valB.Elem()
	}

	// Ensure both values are structs
	if valA.Kind() != reflect.Struct || valB.Kind() != reflect.Struct {
		return
	}

	typA := valA.Type()
	typB := valB.Type()

	// Ensure both structs are of the same type
	if typA != typB {
		return
	}

	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)
		fieldName := typA.Field(i).Name

		// If the field is an embedded struct, compare it recursively
		if fieldA.Kind() == reflect.Struct && fieldB.Kind() == reflect.Struct {
			newParentField := fieldName
			if parentField != "" {
				newParentField = parentField + "." + fieldName
			}
			compareStructs(fieldA, fieldB, newParentField, differences)
		} else {
			// Compare field values
			if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
				fullFieldName := fieldName
				if parentField != "" {
					fullFieldName = parentField + "." + fieldName
				}
				differences = append(differences, fullFieldName)
			}
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Decrypt(encryptionKey string, encryptedString string) (decrypted string, err error) {

	if encryptionKey == "" {
		encryptionKey = os.Getenv("ENCRYPTION_KEY")
	}

	key := []byte(encryptionKey)
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)
	return string(ciphertext), nil
}

func Encrypt(encryptionKey string, plaintextString string) string {
	if encryptionKey == "" {
		encryptionKey = os.Getenv("ENCRYPTION_KEY")
	}
	key := []byte(encryptionKey)
	plaintext := []byte(plaintextString)
	plaintext = PKCS5Padding(plaintext, 16)

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func RemoveItemByValue[T comparable](slice []T, value T) []T {
	newSlice := []T{}
	for _, item := range slice {
		if item != value {
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}

func RemoveItemByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func StructSliceToMapSlice(data interface{}) []map[string]interface{} {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return nil
	}

	sliceLen := dataValue.Len()
	result := make([]map[string]interface{}, sliceLen)

	for i := 0; i < sliceLen; i++ {
		item := dataValue.Index(i)
		if item.Kind() != reflect.Struct {
			return nil
		}

		itemType := item.Type()
		fieldCount := item.NumField()
		itemMap := make(map[string]interface{}, fieldCount)

		for j := 0; j < fieldCount; j++ {
			field := itemType.Field(j)
			fieldValue := item.Field(j)
			itemMap[field.Name] = fieldValue.Interface()
		}

		result[i] = itemMap
	}

	return result
}

func ElementsDiff(src []string, des []string) []string {
	desSet := make(map[string]struct{})
	for _, elem := range des {
		desSet[elem] = struct{}{}
	}

	var result []string
	for _, elem := range src {
		if _, found := desSet[elem]; !found {
			result = append(result, elem)
		}
	}
	return result
}

func GetDay(countryCode string) string {
	timezones := map[string]string{"SA": "Asia/Riyadh", "AE": "Asia/Dubai", "EG": "Africa/Cairo", "FR": "Europe/Paris"}
	timezone, exists := timezones[strings.ToUpper(countryCode)]
	if !exists {
		timezone = "Asia/Riyadh"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Logger.Error(nil, err)
		return strings.ToLower(time.Now().UTC().Weekday().String())
	}
	now := time.Now().In(loc)
	currentDay := now.Weekday()
	day := currentDay.String()
	return strings.ToLower(day)
}

func PrintAsJson(v interface{}) {
	strByte, _ := json.Marshal(v)
	println(string(strByte))
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

func GetAsPointer[T any](p T) *T {
	v := &p
	return v
}

func ArrayToUpper(i []string) []string {
	ii := make([]string, 0)
	if i != nil {
		for _, s := range i {
			ii = append(ii, strings.ToUpper(s))
		}
	}
	return ii
}

func ArrayToLower(i []string) []string {
	ii := make([]string, 0)
	if i != nil {
		for _, s := range i {
			ii = append(ii, strings.ToLower(s))
		}
	}
	return ii
}

func Distance(lt1, lng1, lt2, lng2 float64) float64 {
	lat1 := degreesToRadians(lt1)
	lon1 := degreesToRadians(lng1)
	lat2 := degreesToRadians(lt2)
	lon2 := degreesToRadians(lng2)

	// Earth radius in kilometers
	const earthRadius = 6371.0

	distance := math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*
		math.Cos(lat2)*math.Cos(lon2-lon1)) * earthRadius

	if math.IsNaN(distance) {
		return 0
	}

	return distance
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
func MaskCard(creditCardNumber string) string {
	if len(creditCardNumber) < 15 {
		return "Invalid credit card number"
	}

	lastFourDigits := creditCardNumber[len(creditCardNumber)-4:]

	prefix := creditCardNumber[:len(creditCardNumber)-4]
	padding := strings.Repeat("*", len(prefix))

	maskedCreditCard := padding + lastFourDigits

	return maskedCreditCard
}
func ConvertStructToMap(in interface{}) (response *map[string]interface{}) {
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &response)
	return
}

func IsBearerToken(tokenValue string) (isBearerToken bool, tokenParts []string) {
	tokenParts = strings.Split(tokenValue, " ")
	if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == strings.ToLower("Bearer") {
		isBearerToken = true
	}
	return
}

func GetValueByKey[T any](slice []T, index int) *T {
	if index >= 0 && index < len(slice) {
		return &slice[index]
	}
	return nil
}

func assignMapToStructFields(out interface{}, mapClaims jwt.MapClaims) error {
	outValue := reflect.ValueOf(out).Elem()
	outType := outValue.Type()

	for key, value := range mapClaims {
		found := false
		// Find the corresponding struct field by JSON tag
		for i := 0; i < outType.NumField(); i++ {
			field := outType.Field(i)
			tag := field.Tag.Get("json")
			if tag == key {
				found = true
				fieldValue := outValue.FieldByName(field.Name)
				if fieldValue.IsValid() && fieldValue.CanSet() {
					val := reflect.ValueOf(value)
					if fieldValue.Type() == val.Type() {
						fieldValue.Set(val)
					} else if fieldValue.Kind() == reflect.Int && val.Kind() == reflect.Float64 {
						fieldValue.SetInt(int64(val.Float()))
					} else if fieldValue.Kind() == reflect.String && val.Kind() == reflect.String {
						fieldValue.SetString(val.String())
					}
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("no matching JSON tag found for claim key: %s", key)
		}
	}

	return nil
}

func EqualizeSlices(slice1, slice2 []string) ([]string, []string) {
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})

	for _, s := range slice1 {
		set1[s] = struct{}{}
	}
	for _, s := range slice2 {
		set2[s] = struct{}{}
	}

	var result1, result2 []string

	for _, s := range slice1 {
		if _, exists := set2[s]; exists {
			result1 = append(result1, s)
		}
	}

	for _, s := range slice2 {
		if _, exists := set1[s]; exists {
			result2 = append(result2, s)
		}
	}

	return result1, result2
}
func ObjectToStringified(obj interface{}) string {
	// Marshal the object to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	// Convert JSON bytes to string
	return string(jsonBytes)
}

func SafeMapGet[T any](m map[string]T, key string, defaultValue T) T {
	if m == nil {
		return defaultValue
	}
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

func MarshalUnMarshal(from interface{}, to interface{}) error {
	jsonData, err := json.Marshal(from)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	if err = json.Unmarshal(jsonData, &to); err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	return nil
}

func GetStructName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func CallMethod(obj interface{}, methodName string, args ...interface{}) []reflect.Value {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	method := v.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Println(fmt.Errorf("method not found: %s", methodName))
		return nil
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	return method.Call(in)
}

func TryCatch(f func()) func() error {
	return func() (err error) {
		defer func() {
			if panicInfo := recover(); panicInfo != nil {
				err = fmt.Errorf("%v, %s", panicInfo, string(debug.Stack()))
				return
			}
		}()
		f()
		return err
	}
}

func StructToMap(data interface{}, tagKey string) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		tag := field.Tag.Get(tagKey)
		if tag == "" {
			tag = field.Name
		}
		result[tag] = value
	}
	return result
}

func CopyMapToStruct(doc any, fields map[string]interface{}) error {
	return mapstructure.Decode(fields, doc)
}
