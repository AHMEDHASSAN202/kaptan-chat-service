package utils

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
