package domain

import (
	"context"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/internal/module/menu/responses"
	modifier_group_resp "samm/internal/module/menu/responses/modifier_group"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
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
	GetById(ctx context.Context, id string) (modifier_group_resp.ModifierGroupResp, validators.ErrorResponse)
	List(ctx context.Context, dto *modifier_group.ListModifierGroupsDto) (*responses.ListResponse, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, id string, dto *modifier_group.ChangeModifierGroupStatusDto) validators.ErrorResponse
	SoftDelete(ctx context.Context, id string, input modifier_group.DeleteModifierGroupDto) validators.ErrorResponse
}

type ModifierGroupRepository interface {
	Create(ctx context.Context, docs []ModifierGroup) error
	Update(ctx context.Context, id *primitive.ObjectID, doc *ModifierGroup) error
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]modifier_group_resp.ModifierGroupResp, error)
	List(ctx context.Context, query *modifier_group.ListModifierGroupsDto) ([]ModifierGroup, *mongopagination.PaginationData, error)
	ChangeStatus(ctx context.Context, id *primitive.ObjectID, status *modifier_group.ChangeModifierGroupStatusDto, adminDetails dto.AdminDetails) error
	SoftDelete(ctx context.Context, id *primitive.ObjectID, adminDetails dto.AdminDetails) error
}
