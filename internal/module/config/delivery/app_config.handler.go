package delivery

import (
	"context"
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"
	"samm/pkg/logger"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AppConfigHandler struct {
	appConfigUsecase domain.AppConfigUseCase
	validator        *validator.Validate
	logger           logger.ILogger
}

// InitModifierGroupController will initialize the article's HTTP controller
func InitAppConfigController(e *echo.Echo, appConfigUsecase domain.AppConfigUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &AppConfigHandler{
		appConfigUsecase: appConfigUsecase,
		validator:        validator,
		logger:           logger,
	}
	portal := e.Group("api/v1/admin/app-config")
	{
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.GET("", handler.List)
		portal.GET("/:id", handler.FindById)
		portal.GET("/:type/by-config-type", handler.FindByType)
		portal.DELETE("/:id", handler.Delete)
	}
}

func (a *AppConfigHandler) Create(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input app_config.CreateUpdateAppConfigDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.appConfigUsecase.Create(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *AppConfigHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input app_config.CreateUpdateAppConfigDto
	input.Id = id
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.appConfigUsecase.Update(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *AppConfigHandler) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	config, errResp := a.appConfigUsecase.FindById(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, config)
}

func (a *AppConfigHandler) FindByType(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	configType := c.Param("type")
	if configType == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	config, errResp := a.appConfigUsecase.FindByType(ctx, configType)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, config)
}

func (a *AppConfigHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input app_config.ListAppConfigDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	configs, errResp := a.appConfigUsecase.List(ctx, input)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"configs": configs})
}

func (a *AppConfigHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	errResp := a.appConfigUsecase.SoftDelete(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
