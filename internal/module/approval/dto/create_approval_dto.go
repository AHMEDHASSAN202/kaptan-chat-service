package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/pkg/utils/dto"
)

type CreateApprovalDto struct {
	dto.AdminDetails
	CountryId  string
	EntityId   primitive.ObjectID
	EntityType string
	New        map[string]interface{}
	Old        map[string]interface{}
}
