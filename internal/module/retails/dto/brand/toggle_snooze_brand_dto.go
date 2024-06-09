package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/validators"
)

type BrandToggleSnoozeDto struct {
	Id                    primitive.ObjectID `json:"_" validate:"objectid"`
	IsSnooze              bool               `json:"is_snooze" validate:"required"`
	SnoozeMinutesInterval float32            `json:"snooze_minutes_interval"`
}

func (input *BrandToggleSnoozeDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom validation
	//validate.RegisterValidation("objectid", objectIdValidation)
	//return validators.ValidateStruct(c.Request().Context(), validate, input)
	return validators.ErrorResponse{}
}
