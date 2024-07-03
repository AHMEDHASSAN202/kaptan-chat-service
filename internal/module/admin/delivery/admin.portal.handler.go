package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/admin/consts"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/admin"
	"samm/internal/module/admin/dto/admin_portal"
	"samm/pkg/logger"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/middlewares/portal"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type AdminPortalHandler struct {
	adminUseCase         domain.AdminUseCase
	adminCustomValidator custom_validators.AdminCustomValidator
	validator            *validator.Validate
	logger               logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitAdminPortalController(e *echo.Echo, adminUseCase domain.AdminUseCase, adminCustomValidator custom_validators.AdminCustomValidator, validator *validator.Validate, logger logger.ILogger, portalMiddlewares *portal.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &AdminPortalHandler{
		adminUseCase:         adminUseCase,
		validator:            validator,
		logger:               logger,
		adminCustomValidator: adminCustomValidator,
	}
	portalAdmin := e.Group("api/v1/portal/admin")
	portalAdmin.Use(portalMiddlewares.AuthMiddleware)
	{
		portalAdmin.GET("", handler.ListAdminPortal, commonMiddlewares.PermissionMiddleware("list-portal-admins", "portal-login-accounts"))
		portalAdmin.GET("/:id", handler.FindAdminPortal, commonMiddlewares.PermissionMiddleware("find-portal-admins", "portal-login-accounts"))
		portalAdmin.POST("", handler.CreateAdminPortal, commonMiddlewares.PermissionMiddleware("create-portal-admins", "portal-login-accounts"))
		portalAdmin.PUT("/:id", handler.UpdateAdminPortal, commonMiddlewares.PermissionMiddleware("update-portal-admins", "portal-login-accounts"))
		portalAdmin.DELETE("/:id", handler.DeleteAdminPortal, commonMiddlewares.PermissionMiddleware("delete-portal-admins", "portal-login-accounts"))
		portalAdmin.PUT("/:id/change-status", handler.ChangeStatusAdminPortal, commonMiddlewares.PermissionMiddleware("update-status-portal-admins", "portal-login-accounts"))
	}
}

func (a *AdminPortalHandler) ListAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input admin_portal.ListAdminPortalDTO
	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.adminUseCase.List(ctx, &dto.ListAdminDTO{
		Pagination: input.Pagination,
		Query:      input.Query,
		Status:     input.Status,
		Type:       consts.PORTAL_TYPE,
		Role:       input.Role,
		AccountId:  input.AccountId,
	})
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"admins": data})
}

func (a *AdminPortalHandler) FindAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input dto2.PortalHeaders
	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	data, errResp := a.adminUseCase.Find(ctx, utils.ConvertStringIdToObjectId(id), input.AccountId)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"admin": data})
}

func (a *AdminPortalHandler) CreateAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input admin_portal.CreateAdminPortalDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero(), a.adminCustomValidator.ValidateRoleExists())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.adminUseCase.Create(ctx, &dto.CreateAdminDTO{
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password,
		Type:       consts.PORTAL_TYPE,
		RoleId:     input.RoleId,
		CountryIds: utils.Countries,
		Account:    &dto.Account{Id: utils.ConvertObjectIdToStringId(input.Account.Id), Name: dto.Name{Ar: input.Account.Name.Ar, En: input.Account.Name.En}},
	})
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *AdminPortalHandler) UpdateAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input admin_portal.CreateAdminPortalDTO
	input.ID = utils.ConvertStringIdToObjectId(id)
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero(), a.adminCustomValidator.ValidateRoleExists())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.adminUseCase.Update(ctx, &dto.CreateAdminDTO{
		ID:         input.ID,
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password,
		Type:       consts.PORTAL_TYPE,
		RoleId:     input.RoleId,
		CountryIds: utils.Countries,
		Account:    nil,
		Status:     input.Status,
	})
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *AdminPortalHandler) DeleteAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input admin_portal.DeletePortalAdminDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.adminUseCase.Delete(ctx, utils.ConvertStringIdToObjectId(input.ID), input.AccountId)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *AdminPortalHandler) ChangeStatusAdminPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input admin_portal.ChangeAdminPortalStatusDto
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.adminUseCase.ChangeStatus(ctx, &dto.ChangeAdminStatusDto{
		Id:        input.Id,
		Status:    input.Status,
		AccountId: input.AccountId,
	})
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
