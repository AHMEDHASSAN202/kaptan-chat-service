package collection_method

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type FieldName struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

// basic schema
type CollectionMethodField struct {
	Name       FieldName `json:"name"`
	Type       string    `json:"type"`
	IsRequired bool      `json:"isRequired"`
	Key        string    `json:"key"`
}

type CollectionMethodSchema struct {
	Fields []CollectionMethodField `json:"fields"`
}

// Input payload structure
type Payload map[string]interface{}

func ValidatePayload(ctx context.Context, validate *validator.Validate, schemaMap map[string]interface{}, payload Payload) error {
	schema := CollectionMethodSchema{}
	dbByte, _ := json.Marshal(schemaMap)
	_ = json.Unmarshal(dbByte, &schema)

	for _, field := range schema.Fields {
		value, exists := payload[field.Key]
		fmt.Println(value, exists, field.Key)
		if field.IsRequired && !exists {
			return fmt.Errorf("missing required field: %s", field.Key)
		}

		if exists {
			switch field.Type {
			case "TextPlate", "DropBox":
				if reflect.TypeOf(value).Kind() != reflect.String {
					return fmt.Errorf("field %s should be a string", field.Key)
				}
				if field.IsRequired && value.(string) == "" {
					return fmt.Errorf("field %s should not be empty", field.Key)
				}
			case "CheckBox":
				if reflect.TypeOf(value).Kind() != reflect.Bool {
					return fmt.Errorf("field %s should be a boolean", field.Key)
				}
			default:
				return fmt.Errorf("unsupported field type: %s", field.Type)
			}
		}
	}

	return nil
}
