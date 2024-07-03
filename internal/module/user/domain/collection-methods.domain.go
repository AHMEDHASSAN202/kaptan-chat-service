package domain

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/validators"
	"time"
)

type CollectionMethods struct {
	mgm.DefaultModel `bson:",inline"`
	Type             string             `json:"type" bson:"type"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	Fields           map[string]any     `json:"fields" bson:"fields"`
	Values           map[string]any     `json:"values" bson:"values"`
	DeletedAt        *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type CollectionMethodUseCase interface {
	StoreCollectionMethod(ctx context.Context, payload *CollectionMethods) (err validators.ErrorResponse)
	UpdateCollectionMethod(ctx context.Context, id string, payload *CollectionMethods) (err validators.ErrorResponse)
	FindCollectionMethod(ctx context.Context, Id string, userId string) (user CollectionMethods, err validators.ErrorResponse)
	DeleteCollectionMethod(ctx context.Context, Id string, userId string) (err validators.ErrorResponse)
	ListCollectionMethod(ctx context.Context, collectionMethodType string, userId string) (users []CollectionMethods, err validators.ErrorResponse)
}

type CollectionMethodRepository interface {
	StoreCollectionMethod(ctx context.Context, user *CollectionMethods) (err error)
	UpdateCollectionMethod(ctx context.Context, user *CollectionMethods) (err error)
	FindCollectionMethod(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (user *CollectionMethods, err error)
	DeleteCollectionMethod(ctx context.Context, Id primitive.ObjectID, userId primitive.ObjectID) (err error)
	ListCollectionMethod(ctx context.Context, collectionMethodType string, userId primitive.ObjectID) (locations []CollectionMethods, err error)
}
