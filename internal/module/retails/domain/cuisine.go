package domain

import (
	"context"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/validators"
	"time"
)

type Name struct {
	Ar string `json:"ar" validate:"required,min=3"`
	En string `json:"en" validate:"required,min=3"`
}

type Cuisine struct {
	mgm.DefaultModel `bson:",inline"`
	Name             Name       `json:"name" validate:"required,dive"`
	Logo             string     `json:"logo" bson:"logo"`
	IsHidden         bool       `json:"is_hidden" bson:"is_hidden"`
	DeletedAt        *time.Time `json:"deleted_at" bson:"deleted_at"`
}

type CuisineUseCase interface {
	Create(ctx *context.Context, dto *cuisine.CreateCuisineDto) validators.ErrorResponse
	Update(ctx *context.Context, dto *cuisine.UpdateCuisineDto) validators.ErrorResponse
	Find(ctx *context.Context, id string) (*Cuisine, validators.ErrorResponse)
	GetById(ctx *context.Context, id string) (*Cuisine, validators.ErrorResponse)
	List(ctx *context.Context, dto *cuisine.ListCuisinesDto) (*[]Cuisine, *PaginationData, validators.ErrorResponse)
	ChangeStatus(ctx *context.Context, dto *cuisine.ChangeCuisineStatusDto) validators.ErrorResponse
	SoftDelete(ctx *context.Context, id string) validators.ErrorResponse
	CheckExists(ctx *context.Context, ids []string) validators.ErrorResponse
}

type CuisineRepository interface {
	Create(doc *Cuisine) error
	Update(doc *Cuisine) error
	Find(ctx *context.Context, Id primitive.ObjectID) (*Cuisine, error)
	GetByIds(ctx *context.Context, ids *[]primitive.ObjectID) (*[]Cuisine, error)
	List(ctx *context.Context, query *cuisine.ListCuisinesDto) (*[]Cuisine, *PaginationData, error)
	ChangeStatus(ctx *context.Context, status *cuisine.ChangeCuisineStatusDto) error
	SoftDelete(ctx *context.Context, id primitive.ObjectID) error
}
