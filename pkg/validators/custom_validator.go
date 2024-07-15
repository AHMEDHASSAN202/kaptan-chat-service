package validators

import (
	"context"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"samm/pkg/utils"
	"samm/pkg/validators/localization"
	"time"
)

const TimeFormat = "15:04:05" // Example format (HH:MM:SS)
var (
	saudiRegex = regexp.MustCompile(`^\+9665\d{8}$`)
	egyptRegex = regexp.MustCompile(`^\+201\d{9}$`)
	uaeRegex   = regexp.MustCompile(`^\+9715\d{8}$`)
)

func NewRegisterCustomValidator(c context.Context, validate *validator.Validate) {
	//TODO: context.Background() should depend on the actual context of the request
	registerCustomValidation(c, validate, CustomErrorTags{
		ValidationTag:          localization.Invalid_mongo_ids_validation_rule,
		RegisterValidationFunc: ValidateIDsIsMongoObjectIds,
	}, CustomErrorTags{
		ValidationTag:          localization.Timeformat,
		RegisterValidationFunc: ValidateTimeFormat,
	}, CustomErrorTags{
		ValidationTag:          localization.DateTimeFormat,
		RegisterValidationFunc: ValidateDateTimeFormat,
	}, CustomErrorTags{
		ValidationTag:          localization.PhoneNumber_rule_validation,
		RegisterValidationFunc: PhoneNumberValidator,
	})
}
func ValidateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse(TimeFormat, fl.Field().String())
	return err == nil
}
func ValidateDateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.DateTime, fl.Field().String())
	return err == nil
}

func ValidateIDsIsMongoObjectIds(fl validator.FieldLevel) bool {
	sliceValue := reflect.ValueOf(fl.Field().Interface())
	// Check if the provided slice is actually a slice
	if sliceValue.Kind() == reflect.Slice {
		entityIDs := fl.Field().Interface().([]string)
		if len(entityIDs) == 0 {
			return true
		}
		for _, id := range entityIDs {
			if !utils.IsValidateObjectId(id) {
				return false
			}
		}
	} else if sliceValue.Kind() != reflect.String {
		entityID := fl.Field().Interface().(string)
		if !utils.IsValidateObjectId(entityID) {
			return false
		}
	}

	return true
}

func PhoneNumberValidator(fl validator.FieldLevel) bool {
	// Define a regular expression for validating phone numbers
	countryCode := fl.Parent().FieldByName("CountryCode").String()
	phoneNumber := fl.Field().String()

	fullNumber := countryCode + phoneNumber
	return saudiRegex.MatchString(fullNumber) || egyptRegex.MatchString(fullNumber) || uaeRegex.MatchString(fullNumber)
}
