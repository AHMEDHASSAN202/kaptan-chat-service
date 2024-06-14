package item

import (
	"github.com/jinzhu/copier"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/utils"
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
		itemDocs[i].AdminDetails = append(itemDocs[i].AdminDetails, utils.StructSliceToMapSlice(dto[i].AdminDetails)...)
		itemDocs[i].ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto[i].ModifierGroupsIds)
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
	itemDoc.AdminDetails = append(itemDoc.AdminDetails, utils.StructSliceToMapSlice(dto.AdminDetails)...)
	itemDoc.ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto.ModifierGroupsIds)
}
