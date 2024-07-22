package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/admin"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"time"
)

type AdminHandler struct {
	adminUseCase         domain.AdminUseCase
	adminCustomValidator custom_validators.AdminCustomValidator
	validator            *validator.Validate
	logger               logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitAdminController(e *echo.Echo, adminUseCase domain.AdminUseCase, adminCustomValidator custom_validators.AdminCustomValidator, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &AdminHandler{
		adminUseCase:         adminUseCase,
		validator:            validator,
		logger:               logger,
		adminCustomValidator: adminCustomValidator,
	}
	admin := e.Group("api/v1/admin/admin")
	admin.Use(adminMiddlewares.AuthMiddleware)
	{
		admin.GET("", handler.List, commonMiddlewares.PermissionMiddleware("list-admins"))
		admin.GET("/:id", handler.Find, commonMiddlewares.PermissionMiddleware("find-admins"))
		admin.POST("", handler.Create, commonMiddlewares.PermissionMiddleware("create-admins"))
		admin.PUT("/:id", handler.Update, commonMiddlewares.PermissionMiddleware("update-admins"))
		admin.DELETE("/:id", handler.Delete, commonMiddlewares.PermissionMiddleware("delete-admins"))
		admin.PUT("/:id/change-status", handler.ChangeStatus, commonMiddlewares.PermissionMiddleware("update-status-admins"))
	}
}

func (a *AdminHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.ListAdminDTO
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

	data, errResp := a.adminUseCase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"admins": data})
}

func (a *AdminHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	data, errResp := a.adminUseCase.Find(ctx, utils.ConvertStringIdToObjectId(id), "")
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"admin": data})
}

func (a *AdminHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.CreateAdminDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero(), a.adminCustomValidator.ValidateRoleExists())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.adminUseCase.Create(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *AdminHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input dto.CreateAdminDTO
	input.ID = utils.ConvertStringIdToObjectId(id)

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero(), a.adminCustomValidator.ValidateRoleExists())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.adminUseCase.Update(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *AdminHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	binder := &echo.DefaultBinder{}
	var adminHeaders dto2.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	causerDetails := dto2.AdminDetails{Id: utils.ConvertStringIdToObjectId(adminHeaders.CauserId), Name: adminHeaders.CauserName, Type: adminHeaders.CauserType, Operation: "Delete Admin", UpdatedAt: time.Now()}

	errResp := a.adminUseCase.Delete(ctx, utils.ConvertStringIdToObjectId(id), "", &causerDetails)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *AdminHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.ChangeAdminStatusDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	binder := &echo.DefaultBinder{}
	var adminHeaders dto2.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	causerDetails := dto2.AdminDetails{Id: utils.ConvertStringIdToObjectId(input.CauserId), Name: input.CauserName, Type: input.CauserType, Operation: "Change Admin Status", UpdatedAt: time.Now()}
	input.AdminDetails = causerDetails

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.adminUseCase.ChangeStatus(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
