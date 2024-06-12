package domain

import (
	"context"
	"samm/internal/module/menu/dto/sku"
	"samm/pkg/validators"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SKU struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
}

type SKUUseCase interface {
	Create(ctx context.Context, dto sku.CreateSKUDto) validators.ErrorResponse
	List(ctx context.Context, dto *sku.ListSKUDto) ([]SKU, validators.ErrorResponse)
	CheckExists(ctx context.Context, name string) (bool, validators.ErrorResponse)
}

type SKURepository interface {
	Create(ctx context.Context, doc SKU) error
	List(ctx context.Context, query *sku.ListSKUDto) ([]SKU, error)
	CheckExists(ctx context.Context, name string) (bool, error)
}
