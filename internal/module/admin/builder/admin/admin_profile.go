package admin

import (
	admin2 "samm/internal/module/admin/builder/role"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func AdminProfileBuilder(model *domain.Admin) *admin.AdminProfileResponse {
	adminResp := &admin.AdminProfileResponse{
		ID:         model.ID,
		Name:       model.Name,
		Email:      model.Email,
		Type:       model.Type,
		Role:       admin2.FindRoleBuilder(&model.Role),
		CountryIds: model.CountryIds,
	}
	if model.Account != nil {
		adminResp.Account = &admin.Account{ID: model.Account.Id, Name: admin.LocalizationText{Ar: model.Account.Name.Ar, En: model.Account.Name.En}}
	}
	if model.Kitchen != nil {
		adminResp.Kitchen = &admin.Kitchen{ID: model.Kitchen.Id.Hex(), Name: admin.LocalizationText{Ar: model.Kitchen.Name.Ar, En: model.Kitchen.Name.En}, AllowedStatus: model.Kitchen.AllowedStatus}
	}
	return adminResp
}
