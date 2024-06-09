package menu_group

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	menu_group2 "samm/internal/module/menu/builder/menu_group/dashboard"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/internal/module/menu/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
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
		oRec.logger.Error("MenuGroupUseCase -> Create -> ", err)
		return "", err
	}
	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto)
	menuGroup, errCreate := oRec.repo.Create(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		oRec.logger.Error("MenuGroupUseCase -> Create -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Update(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse) {
	err := oRec.InjectItemsToDTO(ctx, dto)
	if err.IsError {
		oRec.logger.Error("MenuGroupUseCase -> Update -> ", err)
		return "", err
	}
	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto)
	menuGroup, errCreate := oRec.repo.Update(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		oRec.logger.Error("MenuGroupUseCase -> Update -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Delete(ctx context.Context, menuGroupId primitive.ObjectID) validators.ErrorResponse {
	menuGroup, err := oRec.repo.Find(ctx, menuGroupId)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	if menuGroup == nil || menuGroup.ID.IsZero() {
		oRec.logger.Error("MenuGroupUseCase -> Delete -> Error In menuGroup")
		return validators.GetErrorResponse(&ctx, localization.E1002, nil)
	}
	err = oRec.repo.Delete(ctx, menuGroup)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	return validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) ListPortal(ctx context.Context, dto menu_group.ListMenuGroupDTO) (interface{}, validators.ErrorResponse) {
	list, pagination, err := oRec.repo.ListPortal(ctx, dto)
	data := menu_group2.ListGroupBuilder(list)
	listResponse := responses.SetListResponse(data, pagination)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> ListPortal -> ", err)
		return listResponse, validators.GetErrorResponse(&ctx, localization.E1000, nil)
	}
	return listResponse, validators.ErrorResponse{}
}
