package menu_group

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	menu_group2 "samm/internal/module/menu/builder/menu_group/dashboard"
	"samm/internal/module/menu/consts"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/internal/module/menu/external"
	"samm/internal/module/menu/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type MenuGroupUseCase struct {
	repo              domain.MenuGroupRepository
	menuGroupItemRepo domain.MenuGroupItemRepository
	itemRepo          domain.ItemRepository
	logger            logger.ILogger
	extService        external.ExtService
}

func NewMenuGroupUseCase(repo domain.MenuGroupRepository, itemRepo domain.ItemRepository, menuGroupItemRepo domain.MenuGroupItemRepository, logger logger.ILogger, extService external.ExtService) domain.MenuGroupUseCase {
	return &MenuGroupUseCase{
		repo:              repo,
		logger:            logger,
		menuGroupItemRepo: menuGroupItemRepo,
		itemRepo:          itemRepo,
		extService:        extService,
	}
}

func (oRec *MenuGroupUseCase) Create(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse) {
	dto.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Create Menu", UpdatedAt: time.Now()}

	err := oRec.InjectItemsToDTO(ctx, dto)
	if err.IsError {
		oRec.logger.Error("MenuGroupUseCase -> Create -> ", err)
		return "", err
	}

	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto, nil)
	menuGroup, errCreate := oRec.repo.Create(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		oRec.logger.Error("MenuGroupUseCase -> Create -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Update(ctx context.Context, dto *menu_group.CreateMenuGroupDTO) (string, validators.ErrorResponse) {
	menuGroup, errFind := oRec.repo.Find(ctx, dto.ID)
	if errFind != nil {
		oRec.logger.Error("MenuGroupUseCase -> Update -> ", errFind)
		return "", validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	authorized := oRec.AuthorizeMenuGroup(&ctx, menuGroup, utils.ConvertStringIdToObjectId(dto.AccountId))
	if authorized.IsError {
		return "", authorized
	}

	dto.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Update Menu", UpdatedAt: time.Now()}
	err := oRec.InjectItemsToDTO(ctx, dto)
	if err.IsError {
		oRec.logger.Error("MenuGroupUseCase -> Update -> ", err)
		return "", err
	}

	menuGroupDomain, menuGroupItems := menu_group2.MenuGroupBuilder(dto, menuGroup)

	menuGroup, errCreate := oRec.repo.Update(ctx, menuGroupDomain, menuGroupItems)
	if errCreate != nil {
		oRec.logger.Error("MenuGroupUseCase -> Update -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return utils.ConvertObjectIdToStringId(menuGroup.ID), validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Delete(ctx context.Context, dto *menu_group.FindMenuGroupDTO) validators.ErrorResponse {
	menuGroup, err := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(dto.Id))
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	authorized := oRec.AuthorizeMenuGroup(&ctx, menuGroup, utils.ConvertStringIdToObjectId(dto.AccountId))
	if authorized.IsError {
		return authorized
	}

	err = oRec.repo.Delete(ctx, menuGroup)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) List(ctx context.Context, dto *menu_group.ListMenuGroupDTO) (interface{}, validators.ErrorResponse) {
	list, pagination, err := oRec.repo.List(ctx, *dto)
	data := menu_group2.ListGroupBuilder(list)
	listResponse := responses.SetListResponse(data, pagination)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> ListPortal -> ", err)
		return listResponse, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return listResponse, validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) Find(ctx context.Context, dto *menu_group.FindMenuGroupDTO) (interface{}, validators.ErrorResponse) {
	menuGroup, err := oRec.repo.FindWithItems(ctx, utils.ConvertStringIdToObjectId(dto.Id))
	menuGroup = menu_group2.FindMenuGroupBuilder(menuGroup)
	if err != nil {
		oRec.logger.Error("MenuGroupUseCase -> Find -> ", err)
		return menuGroup, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	model := domain.MenuGroup{AccountId: menuGroup.AccountId}
	model.ID = menuGroup.ID
	authorized := oRec.AuthorizeMenuGroup(&ctx, &model, utils.ConvertStringIdToObjectId(dto.AccountId))
	if authorized.IsError {
		return menuGroup, authorized
	}

	branches, errBranches := oRec.extService.RetailsIService.GetBranchesByIds(ctx, utils.ConvertObjectIdsToStringIds(menuGroup.BranchIds))
	if errBranches.IsError {
		oRec.logger.Error("MenuGroupUseCase -> Find -> errBranches -> ", err)
		return menuGroup, validators.GetErrorResponse(&ctx, localization.E1000, nil, utils.GetAsPointer(http.StatusBadRequest))
	}
	menu_group2.PopulateBranches(menuGroup, branches)

	return menuGroup, validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) ChangeStatus(ctx context.Context, id primitive.ObjectID, input *menu_group.ChangeMenuGroupStatusDto) validators.ErrorResponse {
	model, errFind := oRec.repo.Find(ctx, id)
	if errFind != nil {
		oRec.logger.Error("MenuGroupUseCase -> Find -> ChangeStatus -> ", errFind)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	authorized := oRec.AuthorizeMenuGroup(&ctx, model, utils.ConvertStringIdToObjectId(input.AccountId))
	if authorized.IsError {
		return authorized
	}

	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Change Menu Status", UpdatedAt: time.Now()}

	var err error
	if input.Entity == consts.ITEM_CHANGE_STATUS_ENTITY {
		adminDetails.Operation = "Change Menu Item Status"
		err = oRec.menuGroupItemRepo.ChangeItemStatus(ctx, id, input, adminDetails)
	} else if input.Entity == consts.CATEGORY_CHANGE_STATUS_ENTITY {
		adminDetails.Operation = "Change Menu Category Status"
		err = oRec.repo.ChangeCategoryStatus(ctx, model, input, adminDetails)
	} else {
		err = oRec.repo.ChangeMenuStatus(ctx, model, input, adminDetails)
	}

	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
	}

	return validators.ErrorResponse{}
}

func (oRec *MenuGroupUseCase) DeleteEntity(ctx context.Context, input *menu_group.DeleteEntityFromMenuGroupDto) validators.ErrorResponse {
	model, errFind := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.Id))
	if errFind != nil {
		oRec.logger.Error("MenuGroupUseCase -> Find -> DeleteEntity -> ", errFind)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	authorized := oRec.AuthorizeMenuGroup(&ctx, model, utils.ConvertStringIdToObjectId(input.AccountId))
	if authorized.IsError {
		return authorized
	}

	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", UpdatedAt: time.Now()}

	var err error
	if input.Entity == consts.ITEM_CHANGE_STATUS_ENTITY {
		adminDetails.Operation = "Delete Menu Item"
		err = oRec.repo.DeleteItem(ctx, model, input, adminDetails)
	} else if input.Entity == consts.CATEGORY_CHANGE_STATUS_ENTITY {
		adminDetails.Operation = "Delete Menu Category"
		err = oRec.repo.DeleteCategory(ctx, model, input, adminDetails)
	}

	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
	}

	return validators.ErrorResponse{}
}
