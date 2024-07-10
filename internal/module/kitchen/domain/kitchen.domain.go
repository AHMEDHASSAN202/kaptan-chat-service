package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/internal/module/kitchen/responses"
	"samm/pkg/validators"
	"time"
)

type Kitchen struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name       `json:"name" bson:"name"`
	Email            string     `json:"email" bson:"email"`
	Password         string     `json:"-" bson:"password"`
	DeletedAt        *time.Time `json:"-" bson:"deleted_at"`
}

type Name struct {
	Ar string `json:"ar" bson:"ar"`
	En string `json:"en" bson:"en"`
}

type KitchenUseCase interface {
	CreateKitchen(ctx context.Context, payload *kitchen.StoreKitchenDto) (err validators.ErrorResponse)
	UpdateKitchen(ctx context.Context, id string, payload *kitchen.UpdateKitchenDto) (err validators.ErrorResponse)
	FindKitchen(ctx context.Context, Id string) (kitchen Kitchen, err validators.ErrorResponse)
	DeleteKitchen(ctx context.Context, Id string) (err validators.ErrorResponse)
	List(ctx *context.Context, dto *kitchen.ListKitchenDto) (*responses.ListResponse, validators.ErrorResponse)
}

type KitchenRepository interface {
	CreateKitchen(kitchen *Kitchen) (err error)
	UpdateKitchen(kitchen *Kitchen) (err error)
	FindKitchen(ctx context.Context, Id primitive.ObjectID) (kitchen *Kitchen, err error)
	DeleteKitchen(ctx context.Context, Id primitive.ObjectID) (err error)
	List(ctx *context.Context, dto *kitchen.ListKitchenDto) (usersRes *[]Kitchen, paginationMeta *PaginationData, err error)
}

