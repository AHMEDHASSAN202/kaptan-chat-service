package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/validators"
)

type Category struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        LocalizationText     `json:"name" bson:"name"`
	Icon        string               `json:"icon" bson:"icon"`
	Sort        int                  `json:"sort" bson:"sort"`
	Status      string               `json:"status" bson:"status"`
	MenuItemIds []primitive.ObjectID `json:"menu_item_ids" bson:"menu_item_ids"`
}

type MenuGroupAvailability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type MenuGroup struct {
	mgm.DefaultModel `bson:",inline"`
	AccountId        string                  `json:"account_id" bson:"account_id"`
	Name             LocalizationText        `json:"name" bson:"name"`
	BranchIds        []primitive.ObjectID    `json:"branch_ids" bson:"branch_ids"`
	Categories       []Category              `json:"categories" bson:"categories"`
	Availabilities   []MenuGroupAvailability `json:"availabilities" bson:"availabilities"`
	Status           string                  `json:"status" bson:"status"`
}

type MenuGroupUseCase interface {
	Create(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse)
}

type MenuGroupRepository interface {
	Create(ctx context.Context, menuGroup *MenuGroup, menuGroupItems *[]MenuGroupItem) (*MenuGroup, error)
}
