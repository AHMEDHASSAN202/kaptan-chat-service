package role

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type CreateRoleDTO struct {
	ID           primitive.ObjectID `json:"-"`
	Name         Name               `json:"name" validate:"required"`
	Permissions  []string           `json:"permissions" validate:"required"`
	AdminDetails dto.AdminDetails   `json:"-"`
}

func (input *CreateRoleDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
