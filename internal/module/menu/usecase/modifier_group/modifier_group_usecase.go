package modifier_group

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
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
		modifierGroupDocs = append(modifierGroupDocs, convertDtoToCorrespondingDomain(dto[index], nil))

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
	doc := convertDtoToCorrespondingDomain(dto, &oldDoc)
	err := oRec.repo.Update(ctx, &id, &doc)
	if err != nil {
		return validators.GetErrorResponseFromErr(err)
	}
	return validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) GetById(ctx context.Context, id string) (domain.ModifierGroup, validators.ErrorResponse) {
	modifierGroups, err := oRec.repo.GetByIds(ctx, []primitive.ObjectID{utils.ConvertStringIdToObjectId(id)})
	if err != nil {
		return domain.ModifierGroup{}, validators.GetErrorResponseFromErr(err)
	}
	if len(modifierGroups) <= 0 {
		return domain.ModifierGroup{}, validators.GetErrorResponse(&ctx, localization.E1002, nil)
	}
	return modifierGroups[0], validators.ErrorResponse{}
}

func (oRec *ModifierGroupUseCase) List(ctx context.Context, dto *modifier_group.ListModifierGroupsDto) ([]domain.ModifierGroup, validators.ErrorResponse) {
	modifierGroups, err := oRec.repo.List(ctx, dto)
	if err != nil {
		return nil, validators.GetErrorResponseFromErr(err)
	}
	return modifierGroups, validators.ErrorResponse{}
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
