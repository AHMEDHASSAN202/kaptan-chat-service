package modifier_group

import (
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UpdateModifierGroupDto struct {
	Id           string             `json:"_"`
	Name         LocalizationText   `json:"name" validate:"required"`
	Type         string             `json:"type" validate:"required,oneof=required optional"`
	Min          int                `json:"min" validate:"gte=0"`
	Max          int                `json:"max" validate:"gte=0"`
	ProductIds   []string           `json:"product_ids" validate:"product_ids_rules"`
	Status       string             `json:"status" validate:"required,oneof=active inactive"`
	AccountId    string             `json:"account_id" validate:"required,mongodb"`
	AdminDetails []dto.AdminDetails `json:"-"`
	dto.PortalHeaders
}

func (input *UpdateModifierGroupDto) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validate.RegisterValidation("product_ids_rules", utils.ValidateIDsIsMongoObjectIds)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
