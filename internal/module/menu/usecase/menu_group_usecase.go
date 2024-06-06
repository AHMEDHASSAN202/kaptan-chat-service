package usecase

import (
	"context"
	menu_group2 "samm/internal/module/menu/builder/menu_group"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type MenuGroupUseCase struct {
	repo     domain.MenuGroupRepository
	itemRepo domain.ItemRepository
	logger   logger.ILogger
}

func NewMenuGroupUseCase(repo domain.MenuGroupRepository, itemRepo domain.ItemRepository, logger logger.ILogger) domain.MenuGroupUseCase {
	return &MenuGroupUseCase{
		repo:     repo,
		itemRepo: itemRepo,
		logger:   logger,
	}
}

func (oRec *MenuGroupUseCase) Create(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse) {
	err := oRec.InjectItemsToDTO(ctx, dto)
	if err != nil {
		return "", validators.GetErrorResponseFromErr(err)
	}
	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto)
	menuGroup, err := oRec.repo.Create(ctx, menuGroupDomain, menuGroupItems)
	if err != nil {
		return "", validators.GetErrorResponseFromErr(err)
	}
	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}
