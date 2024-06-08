package usecase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err.IsError {
		return "", err
	}
	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto)
	menuGroup, errCreate := oRec.repo.Create(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		return "", validators.GetErrorResponse(&ctx, "E1000", nil)
	}
	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Update(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse) {
	err := oRec.InjectItemsToDTO(ctx, dto)
	if err.IsError {
		return "", err
	}
	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto)
	menuGroup, errCreate := oRec.repo.Update(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		return "", validators.GetErrorResponse(&ctx, "E1000", nil)
	}
	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Delete(ctx context.Context, menuGroupId primitive.ObjectID) validators.ErrorResponse {
	menuGroup, err := oRec.repo.Find(ctx, menuGroupId)
	if err != nil {
		return validators.GetErrorResponse(&ctx, "E1000", nil)
	}
	if menuGroup == nil || menuGroup.ID.IsZero() {
		return validators.GetErrorResponse(&ctx, "E1002", nil)
	}
	err = oRec.repo.Delete(ctx, menuGroup)
	if err != nil {
		return validators.GetErrorResponse(&ctx, "E1000", nil)
	}
	return validators.ErrorResponse{}
}
