package item

import (
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/utils"
	pkgDto "samm/pkg/utils/dto"
	utilsDto "samm/pkg/utils/dto"
	"time"
)

func convertDtoArrToCorrespondingDomain(dto []item.CreateItemDto) []domain.Item {
	itemDocs := make([]domain.Item, 0)
	copier.Copy(&itemDocs, &dto)
	for i, _ := range itemDocs {
		itemDocs[i].AccountId = utils.ConvertStringIdToObjectId(dto[i].AccountId)
		itemDocs[i].DeletedAt = nil
		if dto[i].SKU != "" {
			itemDocs[i].SKU = dto[i].SKU
		}
		itemDocs[i].UpdatedAt = time.Now()
		itemDocs[i].CreatedAt = time.Now()
		if itemDocs[i].AdminDetails == nil {
			itemDocs[i].AdminDetails = make([]pkgDto.AdminDetails, 0)
		}
		itemDocs[i].AdminDetails = append(itemDocs[i].AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto[0].CauserId), Type: dto[0].CauserType, Name: dto[0].CauserName, Operation: "Create", UpdatedAt: time.Now()})

		itemDocs[i].ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto[i].ModifierGroupsIds)
		itemDocs[i].ApprovalStatus = utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL
		if dto[0].CauserType == utils.ADMIN_TYPE {
			itemDocs[i].ApprovalStatus = utils.APPROVAL_STATUS.APPROVED
			itemDocs[i].HasOriginal = true
		}
		itemDocs[i].ID = primitive.NewObjectID()
	}
	return itemDocs
}

func convertDtoToCorrespondingDomain(dto item.UpdateItemDto, itemDoc *domain.Item) {
	copier.Copy(&itemDoc, &dto)
	itemDoc.DeletedAt = nil
	if dto.SKU != "" {
		itemDoc.SKU = dto.SKU
	}
	itemDoc.UpdatedAt = time.Now()
	itemDoc.AccountId = utils.ConvertStringIdToObjectId(dto.AccountId)
	if itemDoc.AdminDetails == nil {
		itemDoc.AdminDetails = make([]pkgDto.AdminDetails, 0)
	}
	itemDoc.AdminDetails = append(itemDoc.AdminDetails, utilsDto.AdminDetails{Id: utils.ConvertStringIdToObjectId(dto.CauserId), Type: dto.CauserType, Name: dto.CauserName, Operation: "Update", UpdatedAt: time.Now()})
	itemDoc.ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto.ModifierGroupsIds)
	itemDoc.ApprovalStatus = utils.APPROVAL_STATUS.WAIT_FOR_APPROVAL
	if dto.CauserType == utils.ADMIN_TYPE {
		itemDoc.ApprovalStatus = utils.APPROVAL_STATUS.APPROVED
		itemDoc.HasOriginal = true
	}
}
