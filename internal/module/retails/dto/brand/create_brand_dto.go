package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3,max=30"`
	En string `json:"en" validate:"required,min=3,max=30"`
}

type CreateBrandDto struct {
	Name        Name     `json:"name" validate:"required"`
	Logo        string   `json:"logo"`
	IsActive    bool     `json:"is_active"`
	SnoozedTill string   `json:"snoozed_till"`
	CuisineIds  []string `json:"cuisine_ids" validate:"cuisine_ids_rule,dive"`
}

func (input *CreateBrandDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
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

func existsIModifierGroup(db interface{}, value []string) bool {
	return true
}
