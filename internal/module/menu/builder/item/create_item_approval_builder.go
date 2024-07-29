package item

import (
	"samm/internal/module/common/dto"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

func CreateItemsApprovalBuilder(items []domain.Item) []dto.CreateApprovalDto {
	approvalData := []dto.CreateApprovalDto{}
	for _, item := range items {
		approvalData = append(approvalData, dto.CreateApprovalDto{
			AdminDetails: item.AdminDetails[len(item.AdminDetails)-1],
			CountryId:    "SA",
			EntityId:     item.ID,
			EntityType:   "items",
			New: map[string]interface{}{
				"name":  utils.StructToMap(item.Name, "bson"),
				"price": item.Price,
				"image": item.Image,
			},
			Old: map[string]interface{}{},
		})
	}
	return approvalData
}

func CreateItemApprovalBuilder(item *domain.Item, n map[string]interface{}, o map[string]interface{}) dto.CreateApprovalDto {
	approvalData := dto.CreateApprovalDto{
		AdminDetails: item.AdminDetails[len(item.AdminDetails)-1],
		CountryId:    "SA",
		EntityId:     item.ID,
		EntityType:   "items",
		New:          n,
		Old:          o,
	}
	return approvalData
}
