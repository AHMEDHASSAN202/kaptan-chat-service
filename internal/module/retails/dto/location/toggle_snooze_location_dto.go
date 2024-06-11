package location

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/validators"
)

type LocationToggleSnoozeDto struct {
	Id                    primitive.ObjectID `json:"_" validate:"objectid"`
	IsSnooze              bool               `json:"is_snooze" validate:"required"`
	SnoozeMinutesInterval float32            `json:"snooze_minutes_interval"`
}

func (input *LocationToggleSnoozeDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ErrorResponse{}
}
