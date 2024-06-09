package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type UpdateBrandDto struct {
	Id string `json:"_"`
	CreateBrandDto
}

func (input *UpdateBrandDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validateModifierGroupsExistsInDB := func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().([]string)
		isValidObjectIds := utils.ValidateIDsIsMongoObjectIds(fl)
		if !isValidObjectIds {
			return false
		}
		return existsIModifierGroup("db", value)
	}
	validate.RegisterValidation("cuisine_ids_rule", validateModifierGroupsExistsInDB)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
