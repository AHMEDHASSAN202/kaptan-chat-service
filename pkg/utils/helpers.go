package utils

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"reflect"
)

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
