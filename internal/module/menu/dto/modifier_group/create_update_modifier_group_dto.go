package modifier_group

import (
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type CreateUpdateModifierGroupDto struct {
	Id           string             `json:"_"`
	Name         LocalizationText   `json:"name" validate:"required,dive"`
	Type         string             `json:"type" validate:"oneof=required optional"`
	Min          int                `json:"min" validate:"required"`
	Max          int                `json:"max" validate:"required"`
	ProductIds   []string           `json:"product_ids" validate:"product_ids_rules"`
	Status       string             `json:"status" validate:"oneof=active inactive"`
	AdminDetails []dto.AdminDetails `json:"-"`
	AccountId    string             `json:"account_id"`
}

func (input *CreateUpdateModifierGroupDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validateProductIdsExistsInDB := func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().([]string)
		isValidObjectIds := utils.ValidateIDsIsMongoObjectIds(fl)
		if !isValidObjectIds {
			return false
		}
		return existsIProduct("db", value)
	}
	validate.RegisterValidation("product_ids_rules", validateProductIdsExistsInDB)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}

func existsIProduct(db interface{}, value []string) bool {
	return true
}
