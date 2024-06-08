package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuGroupItemCategory struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name   LocalizationText   `json:"name" bson:"name"`
	Icon   string             `json:"icon" bson:"icon"`
	Sort   int                `json:"sort" bson:"sort"`
	Status string             `json:"status" bson:"status"`
}

type ItemMenuGroup struct {
	ID             primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	BranchIds      []primitive.ObjectID    `json:"branch_ids" bson:"branch_ids"`
	Availabilities []MenuGroupAvailability `json:"availabilities" bson:"availabilities"`
	Status         string                  `json:"status" bson:"status"`
}

type MenuGroupItem struct {
	mgm.DefaultModel `bson:",inline"`
	ItemId           primitive.ObjectID       `json:"item_id" bson:"item_id"`
	AccountId        string                   `json:"account_id" bson:"account_id"`
	Name             LocalizationText         `json:"name" bson:"name"`
	Desc             LocalizationText         `json:"desc" bson:"desc"`
	Calories         int                      `json:"calories" bson:"calories"`
	Price            float64                  `json:"price" bson:"price"`
	ModifierGroupIds []primitive.ObjectID     `json:"modifier_group_ids" bson:"modifier_group_ids"`
	MenuGroup        ItemMenuGroup            `json:"menu_group" bson:"menu_group"`
	Category         MenuGroupItemCategory    `json:"category" bson:"category"`
	Availabilities   []ItemAvailability       `json:"availabilities" bson:"availabilities"`
	Tags             []string                 `json:"tags" bson:"tags"`
	Image            string                   `json:"image" bson:"image"`
	AdminDetails     []map[string]interface{} `json:"admin_details" bson:"admin_details"`
	Status           string                   `json:"status" bson:"status"`
}

type MenuGroupItemRepository interface {
	CreateUpdateBulk(ctx context.Context, models *[]MenuGroupItem) error
	DeleteBulkByGroupMenuId(ctx context.Context, groupMenuId primitive.ObjectID, exceptionIds []primitive.ObjectID) error
}
