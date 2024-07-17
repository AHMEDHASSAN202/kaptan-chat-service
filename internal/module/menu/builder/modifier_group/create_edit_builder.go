package modifier_group

import (
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"time"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertDtoToCorrespondingDomain(dto modifier_group.CreateModifierGroupDto, causerName string, oldDoc *domain.ModifierGroup) domain.ModifierGroup {
	var modifierGroupDoc domain.ModifierGroup
	copier.Copy(&modifierGroupDoc, &dto)
	modifierGroupDoc.AccountId = utils.ConvertStringIdToObjectId(dto.AccountId)
	modifierGroupDoc.DeletedAt = nil
	modifierGroupDoc.ProductIds = utils.ConvertStringIdsToObjectIds(dto.ProductIds)
	modifierGroupDoc.UpdatedAt = time.Now()
	if oldDoc == nil {
		modifierGroupDoc.CreatedAt = time.Now()
		modifierGroupDoc.AdminDetails = make([]utilsDto.AdminDetails, 0)
		modifierGroupDoc.AdminDetails = append(modifierGroupDoc.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: causerName, Operation: "Create", UpdatedAt: time.Now()})
	} else {
		modifierGroupDoc.CreatedAt = oldDoc.CreatedAt
		modifierGroupDoc.AdminDetails = oldDoc.AdminDetails
		modifierGroupDoc.AdminDetails = append(modifierGroupDoc.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: causerName, Operation: "Update", UpdatedAt: time.Now()})
	}
	return modifierGroupDoc
}

func ConvertUpdateDtoToCorrespondingDomain(dto modifier_group.UpdateModifierGroupDto, causerName string, oldDoc *domain.ModifierGroup) domain.ModifierGroup {
	var modifierGroupDoc domain.ModifierGroup
	copier.Copy(&modifierGroupDoc, &dto)
	modifierGroupDoc.AccountId = oldDoc.AccountId
	modifierGroupDoc.DeletedAt = nil
	modifierGroupDoc.ProductIds = utils.ConvertStringIdsToObjectIds(dto.ProductIds)
	modifierGroupDoc.UpdatedAt = time.Now()

	modifierGroupDoc.CreatedAt = oldDoc.CreatedAt
	modifierGroupDoc.AdminDetails = oldDoc.AdminDetails
	modifierGroupDoc.AdminDetails = append(modifierGroupDoc.AdminDetails, utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: causerName, Operation: "Update", UpdatedAt: time.Now()})

	return modifierGroupDoc
}
