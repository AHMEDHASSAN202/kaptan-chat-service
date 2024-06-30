package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/admin/domain"
	dto "samm/internal/module/admin/dto/auth"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	"samm/pkg/middlewares/portal"
	dto2 "samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type AdminAuthHandler struct {
	adminUseCase domain.AdminUseCase
	validator    *validator.Validate
	logger       logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitAdminAuthController(e *echo.Echo, adminUseCase domain.AdminUseCase, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.Middlewares, portalMiddlewares *portal.Middlewares) {
	handler := &AdminAuthHandler{
		adminUseCase: adminUseCase,
		validator:    validator,
		logger:       logger,
	}

	adminAuth := e.Group("api/v1/admin/auth")
	{
		adminAuth.POST("/login", handler.AdminLogin)
		adminAuth.GET("/profile", handler.AdminProfile, adminMiddlewares.AuthMiddleware)
	}

	portalAuth := e.Group("api/v1/portal/auth")
	{
		portalAuth.POST("/login", handler.PortalLogin)
		portalAuth.GET("/profile", handler.PortalProfile, portalMiddlewares.AuthMiddleware)
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

	profile, errResp := a.adminUseCase.Profile(ctx, input.CauserId)
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

	profile, errResp := a.adminUseCase.Profile(ctx, input.CauserId)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"profile": profile})
}
