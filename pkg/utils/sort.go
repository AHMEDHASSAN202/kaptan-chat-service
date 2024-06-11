package utils

import (
	"fmt"
	"reflect"
	"sort"
)

// SortByField Generic sorting function
func SortByField(slice interface{}, fieldName string) error {
	// Get the slice value and type
	sliceVal := reflect.ValueOf(slice)
	sliceType := reflect.TypeOf(slice)

	// Check if the input is a pointer to a slice
	if sliceType.Kind() != reflect.Ptr || sliceType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("SortByField: expected pointer to a slice, got %v", sliceType)
	}

	// Dereference the pointer to get the actual slice
	sliceVal = sliceVal.Elem()

	// Create a less function using reflection
	lessFunc := func(i, j int) bool {
		itemI := sliceVal.Index(i).FieldByName(fieldName)
		itemJ := sliceVal.Index(j).FieldByName(fieldName)

		// Ensure the field exists and is comparable
		if !itemI.IsValid() || !itemJ.IsValid() {
			return false
		}

		switch itemI.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return itemI.Int() < itemJ.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return itemI.Uint() < itemJ.Uint()
		case reflect.Float32, reflect.Float64:
			return itemI.Float() < itemJ.Float()
		case reflect.String:
			return itemI.String() < itemJ.String()
		default:
			return false
		}
	}

	// Sort the slice using the less function
	sort.Slice(sliceVal.Interface(), lessFunc)

	return nil
}
