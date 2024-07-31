package approval_helper

import (
	"samm/internal/module/approval/dto"
	"samm/internal/module/menu/domain"
	"samm/pkg/utils"
)

func (a *ApprovalItemHelper) CreateItemsApprovalBuilder(items []domain.Item) []dto.CreateApprovalDto {
	approvalData := []dto.CreateApprovalDto{}
	for _, item := range items {
		approvalData = append(approvalData, dto.CreateApprovalDto{
			AdminDetails: item.AdminDetails[len(item.AdminDetails)-1],
			CountryId:    "SA",
			EntityId:     item.ID,
			EntityType:   "items",
			New: map[string]interface{}{
				"name":  utils.StructToMap(item.Name, "bson"),
				"desc":  utils.StructToMap(item.Desc, "bson"),
				"price": item.Price,
				"image": item.Image,
			},
			Old: map[string]interface{}{},
		})
	}
	return approvalData
}

func (a *ApprovalItemHelper) UpdateItemApprovalBuilder(item *domain.Item, n map[string]interface{}, o map[string]interface{}) dto.CreateApprovalDto {
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
