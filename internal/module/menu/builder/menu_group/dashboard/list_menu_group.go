package dashboard

import (
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/responses/menu_group/dashboard"
)

func ListGroupBuilder(models []domain.MenuGroup) []dashboard.ListMenuGroupResponse {
	data := make([]dashboard.ListMenuGroupResponse, 0)
	if models == nil {
		return data
	}
	for _, model := range models {
		data = append(data, dashboard.ListMenuGroupResponse{
			ID:        model.ID,
			AccountId: model.AccountId,
			Name:      model.Name,
			BranchIds: model.BranchIds,
			Status:    model.Status,
			CreatedAt: model.CreatedAt,
			UpdateAt:  model.UpdatedAt,
		})
	}
	return data
}
