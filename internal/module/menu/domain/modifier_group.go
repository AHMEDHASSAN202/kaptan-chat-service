package domain

import (
	"context"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModifierGroup struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name             LocalizationText     `json:"name" bson:"name"`
	Type             string               `json:"type" bson:"type"`
	Min              int                  `json:"min" bson:"min"`
	Max              int                  `json:"max" bson:"max"`
	ProductIds       []primitive.ObjectID `json:"product_ids" bson:"product_ids"`
	Status           string               `json:"status" bson:"status"`
	AdminDetails     []dto.AdminDetails   `json:"admin_details" bson:"admin_details"`
	AccountId        primitive.ObjectID   `json:"account_id" bson:"account_id"`
	DeletedAt        *time.Time           `json:"deleted_at" bson:"deleted_at"`
}

type ModifierGroupUseCase interface {
	Create(ctx context.Context, dto []modifier_group.CreateUpdateModifierGroupDto) validators.ErrorResponse
	Update(ctx context.Context, dto modifier_group.CreateUpdateModifierGroupDto) validators.ErrorResponse
	GetById(ctx context.Context, id string) (ModifierGroup, validators.ErrorResponse)
	List(ctx context.Context, dto *modifier_group.ListModifierGroupsDto) ([]ModifierGroup, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, id string, dto *modifier_group.ChangeModifierGroupStatusDto) validators.ErrorResponse
	SoftDelete(ctx context.Context, id string) validators.ErrorResponse
}

type ModifierGroupRepository interface {
	Create(ctx context.Context, docs []ModifierGroup) error
	Update(ctx context.Context, id *primitive.ObjectID, doc *ModifierGroup) error
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]ModifierGroup, error)
	List(ctx context.Context, query *modifier_group.ListModifierGroupsDto) ([]ModifierGroup, error)
	ChangeStatus(ctx context.Context, id *primitive.ObjectID, status *modifier_group.ChangeModifierGroupStatusDto, adminDetails dto.AdminDetails) error
	SoftDelete(ctx context.Context, id *primitive.ObjectID, adminDetails dto.AdminDetails) error
}
