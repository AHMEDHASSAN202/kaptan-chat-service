package admin

import (
	admin2 "samm/internal/module/admin/builder/role"
	"samm/internal/module/admin/domain"
	"samm/internal/module/admin/responses/admin"
)

func FindAdminBuilder(model *domain.Admin) *admin.FindAdminResponse {
	adminResponse := admin.FindAdminResponse{
		ID:         model.ID,
		Name:       model.Name,
		Email:      model.Email,
		Type:       model.Type,
		Role:       admin2.FindRoleBuilder(&model.Role),
		CountryIds: model.CountryIds,
		Status:     model.Status,
		CreatedAt:  model.CreatedAt,
		UpdateAt:   model.UpdatedAt,
	}
	if model.Account != nil {
		adminResponse.Account = &admin.AccountResp{Id: model.Account.Id, Name: admin.Name{Ar: model.Account.Name.Ar, En: model.Account.Name.En}}
	}
	return &adminResponse
}
