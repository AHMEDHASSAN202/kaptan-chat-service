package menu_group

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/validators"
)

type LocalizationText struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type LocalizationTextDesc struct {
	Ar string `json:"ar"`
	En string `json:"en"`
}

type MenuItemDTO struct {
	Id                string               `json:"id"`
	ItemId            string               `json:"item_id"`
	Sort              int                  `json:"sort"`
	Name              LocalizationText     `json:"name"`
	Desc              LocalizationTextDesc `json:"desc"`
	Type              string               `json:"type"`
	Min               int                  `json:"min"`
	Max               int                  `json:"max"`
	Calories          int                  `json:"calories"`
	Price             float64              `json:"price"`
	ModifierGroupsIds []primitive.ObjectID `json:"modifier_groups_ids"`
	Availabilities    []AvailabilityDTO    `json:"availabilities"`
	Tags              []string             `json:"tags"`
	Image             string               `json:"image"`
	Status            string               `json:"status"`
}

type CategoryDTO struct {
	ID        string           `json:"id"`
	Name      LocalizationText `json:"name"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    string           `json:"status"`
	MenuItems []MenuItemDTO    `json:"menu_items"`
}

type AvailabilityDTO struct {
	Day  string `json:"day" validate:"required_with=From,To"`
	From string `json:"from" validate:"required_with=Day,To"`
	To   string `json:"to"  validate:"required_with=Day,From"`
}

type CreateMenuGroupDTO struct {
	ID             primitive.ObjectID `json:"-"`
	AccountId      string             `json:"account_id"`
	Name           LocalizationText   `json:"name" validate:"required"`
	BranchIds      []string           `json:"branch_ids" validate:"mongodb"`
	Categories     []CategoryDTO      `json:"categories"`
	Availabilities []AvailabilityDTO  `json:"availabilities"`
	Status         string             `json:"status"`
}

func (input *CreateMenuGroupDTO) Validate(c echo.Context, validate *validator.Validate) validators.ErrorResponse {
	return validators.ValidateStruct(c.Request().Context(), validate, input)
}
