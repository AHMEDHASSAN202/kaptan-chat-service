package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/validators"
	"time"
)

type ItemAvailability struct {
	Day  string `json:"day" bson:"day"`
	From string `json:"from" bson:"from"`
	To   string `json:"to" bson:"to"`
}

type Item struct {
	mgm.DefaultModel `bson:",inline"`
	AccountId        primitive.ObjectID       `json:"account_id" bson:"account_id"`
	Name             LocalizationText         `json:"name" bson:"name"`
	Desc             LocalizationText         `json:"desc" bson:"desc"`
	Type             string                   `json:"type" bson:"type"`
	Min              int                      `json:"min" bson:"min"`
	Max              int                      `json:"max" bson:"max"`
	Calories         int                      `json:"calories" bson:"calories"`
	Price            float64                  `json:"price" bson:"price"`
	ModifierGroupIds []primitive.ObjectID     `json:"modifier_groups_ids" bson:"modifier_groups_ids"`
	Availabilities   []ItemAvailability       `json:"availabilities" bson:"availabilities"`
	Tags             []string                 `json:"tags" bson:"tags"`
	Image            string                   `json:"image" bson:"image"`
	AdminDetails     []map[string]interface{} `json:"admin_details" bson:"admin_details"`
	Status           string                   `json:"status" bson:"status"`
	DeletedAt        *time.Time               `json:"deleted_at" bson:"deleted_at"`
}

type ItemUseCase interface {
	Create(ctx context.Context, dto []item.CreateItemDto) validators.ErrorResponse
	Update(ctx context.Context, dto item.UpdateItemDto) validators.ErrorResponse
	GetById(ctx context.Context, id string) (Item, validators.ErrorResponse)
	List(ctx context.Context, dto *item.ListItemsDto) ([]Item, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, id string, dto *item.ChangeItemStatusDto) validators.ErrorResponse
	SoftDelete(ctx context.Context, id string) validators.ErrorResponse
	CheckExists(ctx context.Context, accountId, name string) (bool, validators.ErrorResponse)
}

type ItemRepository interface {
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]Item, error)
	List(ctx context.Context, query *item.ListItemsDto) ([]Item, error)
	Update(ctx context.Context, id *primitive.ObjectID, doc *Item) error
	SoftDelete(ctx context.Context, id *primitive.ObjectID) error
	ChangeStatus(ctx context.Context, id *primitive.ObjectID, status *item.ChangeItemStatusDto) error
	Create(ctx context.Context, doc []Item) error
	CheckExists(ctx context.Context, accountId, name string) (bool, error)
}
