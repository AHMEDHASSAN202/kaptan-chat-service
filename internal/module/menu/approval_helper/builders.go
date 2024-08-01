package approval_helper

import (
	"samm/internal/module/approval/dto"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/external/retails/responses"
	"samm/pkg/utils"
)

func (a *ApprovalItemHelper) CreateItemsApprovalBuilder(items []domain.Item, account responses.AccountByIdResp) []dto.CreateApprovalDto {
	approvalData := []dto.CreateApprovalDto{}
	for _, item := range items {
		approvalData = append(approvalData, dto.CreateApprovalDto{
			AdminDetails: item.AdminDetails[len(item.AdminDetails)-1],
			CountryId:    account.Country.Id,
			EntityId:     item.ID,
			EntityType:   "items",
			New: map[string]interface{}{
				"name":  utils.StructToMap(item.Name, "bson"),
				"desc":  utils.StructToMap(item.Desc, "bson"),
				"price": item.Price,
				"image": item.Image,
			},
			Old:     map[string]interface{}{},
			Doc:     dto.Doc{ID: item.ID, Name: dto.Name{Ar: item.Name.Ar, En: item.Name.En}, Image: item.Image},
			Account: dto.Account{ID: account.ID, Name: dto.Name{Ar: account.Name.Ar, En: account.Name.En}},
		})
	}
	return approvalData
}

func (a *ApprovalItemHelper) UpdateItemApprovalBuilder(item *domain.Item, n map[string]interface{}, o map[string]interface{}, account responses.AccountByIdResp) dto.CreateApprovalDto {
	approvalData := dto.CreateApprovalDto{
		AdminDetails: item.AdminDetails[len(item.AdminDetails)-1],
		CountryId:    account.Country.Id,
		EntityId:     item.ID,
		EntityType:   "items",
		New:          n,
		Old:          o,
		Doc:          dto.Doc{ID: item.ID, Name: dto.Name{Ar: item.Name.Ar, En: item.Name.En}, Image: item.Image},
		Account:      dto.Account{ID: account.ID, Name: dto.Name{Ar: account.Name.Ar, En: account.Name.En}},
	}
	return approvalData
}
