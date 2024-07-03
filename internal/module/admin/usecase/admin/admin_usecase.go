package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	builder "samm/internal/module/admin/builder/admin"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/admin"
	"samm/internal/module/admin/external"
	"samm/internal/module/admin/responses"
	"samm/pkg/database/redis"
	"samm/pkg/jwt"
	"samm/pkg/logger"
	"samm/pkg/utils"
	utilsDto "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type AdminUseCase struct {
	repo             domain.AdminRepository
	roleRepo         domain.RoleRepository
	logger           logger.ILogger
	extService       external.ExtService
	AdminJwtService  jwt.JwtService
	PortalJwtService jwt.JwtService
	redisClient      *redis.RedisClient
}

func NewAdminUseCase(repo domain.AdminRepository, roleRepo domain.RoleRepository, logger logger.ILogger, jwtFactory jwt.JwtServiceFactory, redisClient *redis.RedisClient) domain.AdminUseCase {
	return &AdminUseCase{
		repo:             repo,
		roleRepo:         roleRepo,
		logger:           logger,
		AdminJwtService:  jwtFactory.AdminJwtService(),
		PortalJwtService: jwtFactory.PortalJwtService(),
		redisClient:      redisClient,
	}
}

func (oRec *AdminUseCase) Create(ctx context.Context, input *dto.CreateAdminDTO) (string, validators.ErrorResponse) {
	role, err := oRec.roleRepo.Find(ctx, utils.ConvertStringIdToObjectId(input.RoleId))
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Find Role -> ", err)
		return "", validators.GetErrorResponse(&ctx, localization.E1001, nil, nil)
	}

	input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Create Admin", UpdatedAt: time.Now()}
	adminDomain, err := builder.CreateUpdateAdminBuilder(nil, input, *role)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Create -> ", err)
		return "", validators.GetErrorResponse(&ctx, localization.E1001, nil, nil)
	}

	admin, errCreate := oRec.repo.Create(ctx, adminDomain)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> Create -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return utils.ConvertObjectIdToStringId(admin.ID), validators.ErrorResponse{}
}

func (oRec *AdminUseCase) Update(ctx context.Context, input *dto.CreateAdminDTO) (string, validators.ErrorResponse) {
	admin, errFind := oRec.repo.Find(ctx, input.ID)
	if errFind != nil {
		oRec.logger.Error("AdminUseCase -> Update -> ", errFind)
		return "", validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	if input.Account != nil && input.Account.Id != "" && !admin.Authorized(input.Account.Id) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Update Admin -> ", admin.ID)
		return "", validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	role, err := oRec.roleRepo.Find(ctx, utils.ConvertStringIdToObjectId(input.RoleId))
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Find Role -> ", err)
		return "", validators.GetErrorResponse(&ctx, localization.E1001, nil, nil)
	}

	input.AdminDetails = utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Update Admin", UpdatedAt: time.Now()}
	adminDomain, err := builder.CreateUpdateAdminBuilder(admin, input, *role)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Update -> ", err)
	}

	admin, errCreate := oRec.repo.Update(ctx, adminDomain)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> Update -> ", errCreate)
		return "", validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	oRec.RemoveAdminFromCache(utils.ConvertObjectIdToStringId(admin.ID))

	return utils.ConvertObjectIdToStringId(admin.ID), validators.ErrorResponse{}
}

func (oRec *AdminUseCase) Delete(ctx context.Context, adminId primitive.ObjectID, accountId string) validators.ErrorResponse {
	admin, err := oRec.repo.Find(ctx, adminId)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	if accountId != "" && !admin.Authorized(accountId) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Update Admin -> ", admin.ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Delete Admin", UpdatedAt: time.Now()}
	err = oRec.repo.Delete(ctx, admin, adminDetails)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Delete -> ", err)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}

	oRec.RemoveAdminFromCache(utils.ConvertObjectIdToStringId(admin.ID))

	return validators.ErrorResponse{}
}

func (oRec *AdminUseCase) List(ctx context.Context, input *dto.ListAdminDTO) (interface{}, validators.ErrorResponse) {
	list, pagination, err := oRec.repo.List(ctx, input)
	data := builder.ListAdminBuilder(&list)
	listResponse := responses.SetListResponse(data, pagination)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> List -> ", err)
		return listResponse, validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return listResponse, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) Find(ctx context.Context, adminId primitive.ObjectID, accountId string) (interface{}, validators.ErrorResponse) {
	admin, err := oRec.repo.Find(ctx, adminId)
	adminResp := builder.FindAdminBuilder(admin)
	if err != nil {
		oRec.logger.Error("AdminUseCase -> Find -> ", err)
		return adminResp, validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	if accountId != "" && !admin.Authorized(accountId) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized Find Admin -> ", admin.ID)
		return nil, validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	return adminResp, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) ChangeStatus(ctx context.Context, input *dto.ChangeAdminStatusDto) validators.ErrorResponse {
	admin, errFind := oRec.repo.Find(ctx, utils.ConvertStringIdToObjectId(input.Id))
	if errFind != nil {
		oRec.logger.Error("AdminUseCase -> Find -> ChangeStatus -> ", errFind)
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, utils.GetAsPointer(http.StatusNotFound))
	}

	if input.AccountId != "" && !admin.Authorized(input.AccountId) {
		oRec.logger.Error("AuthorizeMenuGroup -> UnAuthorized ChangeStatus Admin -> ", admin.ID)
		return validators.GetErrorResponse(&ctx, localization.E1006, nil, utils.GetAsPointer(http.StatusForbidden))
	}

	adminDetails := utilsDto.AdminDetails{Id: primitive.NewObjectID(), Name: "Hassan", Operation: "Change Admin Status", UpdatedAt: time.Now()}
	err := oRec.repo.ChangeStatus(ctx, admin, input, adminDetails)
	if err != nil {
		return validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)
	}

	oRec.RemoveAdminFromCache(utils.ConvertObjectIdToStringId(admin.ID))

	return validators.ErrorResponse{}
}

func (oRec *AdminUseCase) CheckEmailExists(ctx context.Context, name string, adminId primitive.ObjectID) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckEmailExists(ctx, name, adminId)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) CheckRoleExists(ctx context.Context, roleId primitive.ObjectID) (bool, validators.ErrorResponse) {
	isExists, err := oRec.repo.CheckRoleExists(ctx, roleId)
	if err != nil {
		return isExists, validators.GetErrorResponseFromErr(err)
	}
	return isExists, validators.ErrorResponse{}
}

func (oRec *AdminUseCase) SyncAccount(ctx context.Context, input dto.Account) validators.ErrorResponse {
	errCreate := oRec.repo.SyncAccount(ctx, input)
	if errCreate != nil {
		oRec.logger.Error("AdminUseCase -> UpdateAccountById -> ", errCreate)
		return validators.GetErrorResponse(&ctx, localization.E1000, nil, nil)
	}
	return validators.ErrorResponse{}
}
