package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/brand"
	"samm/internal/module/retails/responses"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type Brand struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name                 `json:"name" validate:"required"`
	Logo             string               `json:"logo" bson:"logo"`
	IsActive         bool                 `json:"is_active" bson:"is_active"`
	CuisineIds       []primitive.ObjectID `json:"cuisine_ids" validate:"cuisine_ids_rule"`
	DeletedAt        *time.Time           `json:"deleted_at" bson:"deleted_at"`
	Cuisines         *[]Cuisine           `json:"cuisines" bson:"cuisines"`
	AdminDetails     []dto.AdminDetails   `json:"admin_details" bson:"admin_details,omitempty"`
}

type BrandUseCase interface {
	Create(ctx context.Context, dto *brand.CreateBrandDto) (*Brand, validators.ErrorResponse)
	Update(ctx *context.Context, dto *brand.UpdateBrandDto) validators.ErrorResponse
	Find(ctx *context.Context, id string) (*Brand, validators.ErrorResponse)
	FindWithCuisines(ctx *context.Context, id string) (*Brand, validators.ErrorResponse)
	GetById(ctx *context.Context, id string) (*Brand, validators.ErrorResponse)
	List(ctx *context.Context, dto *brand.ListBrandDto) (*responses.ListResponse, validators.ErrorResponse)
	ChangeStatus(ctx *context.Context, dto *brand.ChangeBrandStatusDto) validators.ErrorResponse
	SoftDelete(ctx *context.Context, id string, adminDetails *dto.AdminHeaders) validators.ErrorResponse
}

type BrandRepository interface {
	Create(ctx context.Context, doc *Brand) error
	Update(doc *Brand) error
	FindBrand(*context.Context, primitive.ObjectID) (*Brand, error)
	FindWithCuisines(context.Context, primitive.ObjectID) (*Brand, error)
	GetByIds(ctx *context.Context, ids *[]primitive.ObjectID) (*[]Brand, error)
	List(ctx *context.Context, query *brand.ListBrandDto) (*[]Brand, *PaginationData, error)
	UpdateBrandAndLocations(doc *Brand) error
	SoftDelete(doc *Brand) error
}

//func (model *Brand) Updated(ctx context.Context, result *mongo.UpdateResult) error {

//	return nil
//}
