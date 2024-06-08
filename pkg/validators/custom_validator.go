package validators

import (
	"github.com/go-playground/validator/v10"
	"time"
)

const TimeFormat = "15:04:05" // Example format (HH:MM:SS)

func ValidateTimeFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse(TimeFormat, fl.Field().String())
	return err == nil
}
