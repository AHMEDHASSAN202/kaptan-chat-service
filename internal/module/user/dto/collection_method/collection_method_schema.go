package collection_method

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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

func ValidatePayload(ctx context.Context, validate *validator.Validate, schemaMap map[string]interface{}, payload Payload) validators.ErrorResponse {
	schema := CollectionMethodSchema{}
	dbByte, _ := json.Marshal(schemaMap)
	_ = json.Unmarshal(dbByte, &schema)

	for _, field := range schema.Fields {
		value, exists := payload[field.Key]
		if field.IsRequired && !exists {
			fmt.Println(fmt.Sprintf("missing required field: %s", field.Key))
			return validators.GetErrorResponse(&ctx, localization.E1008, nil, nil)
		}

		if exists {
			switch field.Type {
			case "TextPlate", "Text", "DropBox":
				if reflect.TypeOf(value).Kind() != reflect.String {
					fmt.Println(fmt.Sprintf("field %s should be a string", field.Key))
					return validators.GetErrorResponse(&ctx, localization.E1008, nil, nil)
				}
				if field.IsRequired && value.(string) == "" {
					fmt.Println(fmt.Sprintf("field %s should not be empty", field.Key))
					return validators.GetErrorResponse(&ctx, localization.E1008, nil, nil)
				}
			case "CheckBox":
				if reflect.TypeOf(value).Kind() != reflect.Bool {
					fmt.Println(fmt.Sprintf("field %s should be a boolean", field.Key))
					return validators.GetErrorResponse(&ctx, localization.E1008, nil, nil)
				}
			default:
				fmt.Println(fmt.Sprintf("unsupported field type: %s", field.Type))
				return validators.GetErrorResponse(&ctx, localization.E1008, nil, nil)
			}
		}
	}

	return validators.ErrorResponse{}
}
