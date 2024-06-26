package brand

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type UpdateBrandDto struct {
	Id string `json:"_"`
	CreateBrandDto
}

func (input *UpdateBrandDto) Validate(c echo.Context, validate *validator.Validate, validateCuisineExists func(fl validator.FieldLevel) bool) validators.ErrorResponse {
	validateModifierGroupsExistsInDB := func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().([]string)
		isValidObjectIds := utils.ValidateIDsIsMongoObjectIds(fl)
		if !isValidObjectIds {
			return false
		}
		return existsIModifierGroup("db", value)
	}
	validate.RegisterValidation("cuisine_ids_rule", validateModifierGroupsExistsInDB)
	return validators.ValidateStruct(c.Request().Context(), validate, input, validators.CustomErrorTags{
		ValidationTag:          localization.Cuisine_id_is_exists_rules_validation,
		RegisterValidationFunc: validateCuisineExists,
	})
}
