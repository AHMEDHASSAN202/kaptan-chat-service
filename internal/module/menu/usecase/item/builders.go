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
		itemDocs[i].UpdatedAt = time.Now()
		itemDocs[i].CreatedAt = time.Now()
		itemDocs[i].ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto[i].ModifierGroupsIds)

	}
	return itemDocs
}
func convertDtoToCorrespondingDomain(dto item.UpdateItemDto) domain.Item {
	itemDoc := domain.Item{}
	copier.Copy(&itemDoc, &dto)
	itemDoc.DeletedAt = nil
	itemDoc.UpdatedAt = time.Now()
	itemDoc.ModifierGroupIds = utils.ConvertStringIdsToObjectIds(dto.ModifierGroupsIds)
	return itemDoc
}
