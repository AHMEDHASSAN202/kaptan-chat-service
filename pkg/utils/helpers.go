package utils

import (
	"encoding/json"
	"fmt"
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

// LogJSON logs data as JSON
func LogJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		return
	}
	fmt.Println("Logged JSON:")
	fmt.Println(string(jsonData))
}
