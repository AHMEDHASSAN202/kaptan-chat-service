package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type LocalizationTextWithoutValidation struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type MenuItemDTO struct {
	Id               string                            `json:"id"`
	ItemId           string                            `json:"item_id" validate:"required,mongodb"`
	Sort             int                               `json:"sort"`
	Name             LocalizationTextWithoutValidation `json:"name"`
	Desc             LocalizationTextWithoutValidation `json:"desc"`
	Type             string                            `json:"type"`
	Min              int                               `json:"min"`
	Max              int                               `json:"max"`
	Calories         int                               `json:"calories"`
	Price            float64                           `json:"price"`
	ModifierGroupIds []primitive.ObjectID              `json:"modifier_group_ids"`
	Availabilities   []AvailabilityDTO                 `json:"availabilities"`
	Tags             []string                          `json:"tags"`
	Image            string                            `json:"image"`
	Status           string                            `json:"status" validate:"oneof=active inactive"`
}

type CategoryDTO struct {
	ID        string           `json:"id"`
	Name      LocalizationText `json:"name" validate:"required"`
	Icon      string           `json:"icon" validate:"required,url"`
	Sort      int              `json:"sort" validate:"required"`
	Status    string           `json:"status" validate:"oneof=active inactive"`
	MenuItems []MenuItemDTO    `json:"menu_items" validate:"dive"`
}

type AvailabilityDTO struct {
	Day  string `json:"day" validate:"required_with=From"`
	From string `json:"from" validate:"required_with=Day"`
	To   string `json:"to"  validate:"required_with=From"`
}

type CreateMenuGroupDTO struct {
	ID             primitive.ObjectID `json:"-"`
	AccountId      string             `json:"account_id" validate:"mongodb"`
	Name           LocalizationText   `json:"name" validate:"required"`
	BranchIds      []string           `json:"branch_ids" validate:"branch_ids_rules"`
	Categories     []CategoryDTO      `json:"categories" validate:"dive"`
	Availabilities []AvailabilityDTO  `json:"availabilities" validate:"dive"`
	Status         string             `json:"status" validate:"oneof=active inactive"`
}

func (input *CreateMenuGroupDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	validate.RegisterValidation("branch_ids_rules", utils.ValidateIDsIsMongoObjectIds)
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
