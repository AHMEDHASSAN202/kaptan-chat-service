package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/admin/custom_validators"
	"samm/internal/module/admin/domain"
	admin2 "samm/internal/module/admin/dto/admin"
	dto "samm/internal/module/admin/dto/auth"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/middlewares/portal"
	"samm/pkg/utils"
	dto2 "samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type AdminAuthHandler struct {
	adminUseCase         domain.AdminUseCase
	validator            *validator.Validate
	logger               logger.ILogger
	adminCustomValidator custom_validators.AdminCustomValidator
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitAdminAuthController(e *echo.Echo, adminUseCase domain.AdminUseCase, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares, portalMiddlewares *portal.ProviderMiddlewares, adminCustomValidator custom_validators.AdminCustomValidator, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &AdminAuthHandler{
		adminUseCase:         adminUseCase,
		validator:            validator,
		logger:               logger,
		adminCustomValidator: adminCustomValidator,
	}

	adminAuth := e.Group("api/v1/admin/auth")
	{
		adminAuth.POST("/login", handler.AdminLogin)
		adminAuth.GET("/profile", handler.AdminProfile, adminMiddlewares.AuthMiddleware)
		adminAuth.PUT("/profile", handler.UpdateAdminProfile, adminMiddlewares.AuthMiddleware)
		adminAuth.POST("/login-as-portal/:id", handler.LoginAsPortal, adminMiddlewares.AuthMiddleware, commonMiddlewares.PermissionMiddleware("portal-login-accounts"))
	}

	portalAuth := e.Group("api/v1/portal/auth")
	{
		portalAuth.POST("/login", handler.PortalLogin)
		portalAuth.GET("/profile", handler.PortalProfile, portalMiddlewares.AuthMiddleware)
		portalAuth.PUT("/profile", handler.UpdatePortalProfile, portalMiddlewares.AuthMiddleware)
	}
}

func (a *AdminAuthHandler) AdminLogin(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.AdminAuthDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	profile, token, errResp := a.adminUseCase.AdminLogin(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile, "token": token})
}

func (a *AdminAuthHandler) PortalLogin(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.PortalAuthDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	profile, token, errResp := a.adminUseCase.PortalLogin(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile, "token": token})
}

func (a *AdminAuthHandler) AdminProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto2.AdminHeaders
	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	profile, errResp := a.adminUseCase.Profile(ctx, dto.ProfileDTO{AdminId: input.CauserId, AccountId: "", CauserDetails: nil})
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile})
}

func (a *AdminAuthHandler) PortalProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto2.PortalHeaders
	binder := &echo.DefaultBinder{}
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	profile, errResp := a.adminUseCase.Profile(ctx, dto.ProfileDTO{AdminId: input.CauserId, AccountId: input.CauserAccountId, CauserDetails: input.GetCauserDetailsAsMap()})
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile})
}

func (a *AdminAuthHandler) UpdateAdminProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.UpdateAdminProfileDTO
	binder := &echo.DefaultBinder{}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	input.ID = utils.ConvertStringIdToObjectId(input.CauserId)
	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	profile, errResp := a.adminUseCase.UpdateAdminProfile(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile})
}

func (a *AdminAuthHandler) UpdatePortalProfile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input dto.UpdatePortalProfileDTO
	binder := &echo.DefaultBinder{}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	input.ID = utils.ConvertStringIdToObjectId(input.CauserId)
	validationErr := input.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.adminCustomValidator.PasswordRequiredIfIdIsZero())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	profile, errResp := a.adminUseCase.UpdatePortalProfile(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile})
}

func (a *AdminAuthHandler) LoginAsPortal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input admin2.LoginAsPortalDto
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
		return validators.ErrorStatusBadRequest(c, validationErr)
	}

	profile, token, errResp := a.adminUseCase.LoginAsPortal(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile, "token": token})
}
