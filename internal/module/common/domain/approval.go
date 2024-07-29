package domain

import (
	"context"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	dto2 "samm/internal/module/common/dto"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"time"
)

type Dates struct {
	ApprovedAt *time.Time `json:"approved_at" bson:"approved_at"`
	RejectedAt *time.Time `json:"rejected_at" bson:"rejected_at"`
}

type Fields struct {
	New map[string]interface{} `json:"new" bson:"new"`
	Old map[string]interface{} `json:"old" bson:"old"`
}

type Approval struct {
	mgm.DefaultModel `bson:",inline"`
	CountryId        string             `json:"country_id" bson:"country_id"`
	EntityId         primitive.ObjectID `json:"entity_id" bson:"entity_id"`
	EntityType       string             `json:"entity_type" bson:"entity_type"`
	Fields           Fields             `json:"fields" bson:"fields"`
	Status           string             `json:"status" bson:"status"`
	Dates            Dates              `json:"dates" bson:"dates"`
	AdminDetails     dto.AdminDetails   `json:"admin_details" bson:"admin_details,omitempty"`
}

type ApprovalUseCase interface {
	List(ctx context.Context, dto *dto2.ListApprovalDto) (interface{}, validators.ErrorResponse)
	FindByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) (interface{}, validators.ErrorResponse)
	ChangeStatus(ctx context.Context, approvalId *dto2.ChangeStatusApprovalDto) validators.ErrorResponse
}

type ApprovalRepository interface {
	CreateOrUpdate(ctx context.Context, dto []dto2.CreateApprovalDto) error
	FindByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) (*Approval, error)
	FindById(ctx context.Context, approvalId primitive.ObjectID) (*Approval, error)
	List(ctx context.Context, dto *dto2.ListApprovalDto) ([]Approval, *mongopagination.PaginationData, error)
	ChangeStatus(ctx context.Context, domainData *Approval) error
	DeleteByEntity(ctx context.Context, entityId primitive.ObjectID, entityType string) error
}
