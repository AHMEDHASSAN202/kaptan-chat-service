package role

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	builder "samm/internal/module/admin/builder/role"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/role"
	"samm/internal/module/admin/responses"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type RoleUseCase struct {
	repo   domain.RoleRepository
	logger logger.ILogger
}

func NewRoleUseCase(repo domain.RoleRepository, logger logger.ILogger) domain.RoleUseCase {
	return &RoleUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (oRec *RoleUseCase) Create(ctx context.Context, input *dto.CreateRoleDTO) (string, validators.ErrorResponse) {
	input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Create Role", UpdatedAt: time.Now()}
	roleDomain, err := builder.CreateUpdateRoleBuilder(nil, input)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> Create -> ", err)
		return "", validators.GetErrorResponse(&ctx, localization.E1001, nil, nil)
	}

	role, errCreate := oRec.repo.Create(ctx, roleDomain)
	if errCreate != nil {
		oRec.logger.Error("RoleUseCase -> Create -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return utils.ConvertObjectIdToStringId(role.ID), validators.ErrorResponse{}
}

func (oRec *RoleUseCase) Update(ctx context.Context, input *dto.CreateRoleDTO) (string, validators.ErrorResponse) {
	role, errFind := oRec.repo.Find(ctx, input.ID)
	if errFind != nil {
		oRec.logger.Error("RoleUseCase -> Update -> ", errFind)
		return "", validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Update Role", UpdatedAt: time.Now()}
	roleDomain, err := builder.CreateUpdateRoleBuilder(role, input)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> Update -> ", err)
	}

	role, errCreate := oRec.repo.Update(ctx, roleDomain)
	if errCreate != nil {
		oRec.logger.Error("RoleUseCase -> Update -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return utils.ConvertObjectIdToStringId(role.ID), validators.ErrorResponse{}
}

func (oRec *RoleUseCase) Delete(ctx context.Context, roleId primitive.ObjectID) validators.ErrorResponse {
	role, err := oRec.repo.Find(ctx, roleId)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	err = oRec.repo.Delete(ctx, role)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	return validators.ErrorResponse{}
}

func (oRec *RoleUseCase) List(ctx context.Context, input *dto.ListRoleDTO) (interface{}, validators.ErrorResponse) {
	list, pagination, err := oRec.repo.List(ctx, input)
	data := builder.ListRoleBuilder(&list)
	listResponse := responses.SetListResponse(data, pagination)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> List -> ", err)
		return listResponse, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return listResponse, validators.ErrorResponse{}
}

func (oRec *RoleUseCase) Find(ctx context.Context, roleId primitive.ObjectID) (interface{}, validators.ErrorResponse) {
	role, err := oRec.repo.Find(ctx, roleId)
	roleResp := builder.FindRoleBuilder(role)
	if err != nil {
		oRec.logger.Error("RoleUseCase -> Find -> ", err)
		return roleResp, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}
	return roleResp, validators.ErrorResponse{}
}
