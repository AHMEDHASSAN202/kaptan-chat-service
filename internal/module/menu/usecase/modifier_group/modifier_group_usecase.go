package modifier_group

import (
	"context"
	"github.com/jinzhu/copier"
	builder "samm/internal/module/menu/builder/modifier_group"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/internal/module/menu/responses"
	modifier_group_resp "samm/internal/module/menu/responses/modifier_group"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModifierGroupUseCase struct {
	repo   domain.ModifierGroupRepository
	logger logger.ILogger
}

func NewModifierGroupUseCase(repo domain.ModifierGroupRepository, logger logger.ILogger) domain.ModifierGroupUseCase {
	return &ModifierGroupUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *ModifierGroupUseCase) Create(ctx context.Context, dto []modifier_group.CreateUpdateModifierGroupDto) validators.ErrorResponse {

	modifierGroupDocs := make([]domain.ModifierGroup, 0)
	for index := range dto {
		modifierGroupDocs = append(modifierGroupDocs, builder.ConvertDtoToCorrespondingDomain(dto[index], nil))

	}

	err := oRec.repo.Create(ctx, modifierGroupDocs)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) Update(ctx context.Context, dto modifier_group.CreateUpdateModifierGroupDto) validators.ErrorResponse {
	oldDoc, getByIdErr := oRec.GetById(ctx, dto.Id)
	if getByIdErr.IsError {
		return getByIdErr
	}
	id := utils.ConvertStringIdToObjectId(dto.Id)
	modifierGroup := domain.ModifierGroup{}
	copier.Copy(&modifierGroup, &oldDoc)
	doc := builder.ConvertDtoToCorrespondingDomain(dto, &modifierGroup)
	err := oRec.repo.Update(ctx, &id, &doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) GetById(ctx context.Context, id string) (modifier_group_resp.ModifierGroupResp, validators.ErrorResponse) {
	modifierGroups, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return modifier_group_resp.ModifierGroupResp{}, validators.GetErrorResponseFromErr(err)
	}
	if len(modifierGroups) <= 0 {
		return modifier_group_resp.ModifierGroupResp{}, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
	}
	modifierGroupsResp := modifierGroups[0]
	if modifierGroupsResp.Products == nil {
		modifierGroupsResp.Products = make([]map[string]any, 0)
	}
	return modifierGroupsResp, validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) List(ctx context.Context, dto *modifier_group.ListModifierGroupsDto) (*responses.ListResponse, validators.ErrorResponse) {
	modifierGroups, paginationResult, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}

	return responses.SetListResponse(modifierGroups, paginationResult), validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) ChangeStatus(ctx context.Context, id string, dto *modifier_group.ChangeModifierGroupStatusDto) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Change Status", UpdatedAt: time.Now()}
	err := oRec.repo.ChangeStatus(ctx, &idDoc, dto, adminDetails)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) SoftDelete(ctx context.Context, id string) validators.ErrorResponse {
	idDoc := utils.ConvertStringIdToObjectId(id)
	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Malhat", Operation: "Delete", UpdatedAt: time.Now()}
	err := oRec.repo.SoftDelete(ctx, &idDoc, adminDetails)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}
