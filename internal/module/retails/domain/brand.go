package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/validators"
	"time"
)

type Brand struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name                 `json:"name" validate:"required"`
	Logo             string               `json:"logo" bson:"logo"`
	IsActive         bool                 `json:"is_active" bson:"is_active"`
	SnoozedTill      string               `json:"snoozed_till" bson:"snoozed_till"`
	CuisineIds       []primitive.ObjectID `json:"cuisine_ids" validate:"cuisine_ids_rule"`
	DeletedAt        *time.Time           `json:"deleted_at" bson:"deleted_at"`
}

type BrandUseCase interface {
	Create(ctx *context.Context, dto *brand.CreateBrandDto) validators.ErrorResponse
	Update(ctx *context.Context, dto *brand.UpdateBrandDto) validators.ErrorResponse
	Find(ctx *context.Context, id string) (*Brand, validators.ErrorResponse)
	GetById(ctx *context.Context, id string) (*Brand, validators.ErrorResponse)
	List(ctx *context.Context, dto *brand.ListBrandDto) (*[]Brand, *PaginationData, validators.ErrorResponse)
	ChangeStatus(ctx *context.Context, dto *brand.ChangeBrandStatusDto) validators.ErrorResponse
	ToggleSnooze(ctx *context.Context, dto *brand.BrandToggleSnoozeDto) validators.ErrorResponse
	SoftDelete(ctx *context.Context, id string) validators.ErrorResponse
}

type BrandRepository interface {
	Create(doc *Brand) error
	Update(doc *Brand) error
	FindBrand(*context.Context, primitive.ObjectID) (*Brand, error)
	GetByIds(ctx *context.Context, ids *[]primitive.ObjectID) (*[]Brand, error)
	List(ctx *context.Context, query *brand.ListBrandDto) (*[]Brand, *PaginationData, error)
	UpdateBrandAndLocations(doc *Brand) error
	SoftDelete(doc *Brand) error
}
