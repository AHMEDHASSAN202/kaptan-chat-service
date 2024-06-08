package cuisine

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/validators"
)

type ChangeCuisineStatusDto struct {
	Id       primitive.ObjectID `json:"_" validate:"objectid"`
	IsHidden bool               `json:"is_hidden" validate:"required"`
}

func (input *ChangeCuisineStatusDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	// Register custom validation
	//validate.RegisterValidation("objectid", objectIdValidation)
	//return validators.ValidateStruct(c.Request().Context(), validate, input)
	return validators.ErrorResponse{}
}

// objectIdValidation checks if a primitive.ObjectID is valid (not the zero value).
func objectIdValidation(fl validator.FieldLevel) bool {
	id := fl.Field().Interface().(primitive.ObjectID)
	return id != primitive.NilObjectID
}
