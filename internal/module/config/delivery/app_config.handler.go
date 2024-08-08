package delivery

import (
	"context"
	"samm/internal/module/config/consts"
	"samm/internal/module/config/custom_validators"
	"samm/internal/module/config/domain"
	"samm/internal/module/config/dto/app_config"
	"samm/pkg/app_localization"
	"samm/pkg/logger"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AppConfigHandler struct {
	appConfigUsecase         domain.AppConfigUseCase
	appConfigCustomValidator custom_validators.AppConfigCustomValidator
	validator                *validator.Validate
	logger                   logger.ILogger
	localizationData         map[string]interface{}
}

// InitAppConfigController will initialize the article's HTTP controller
func InitAppConfigController(e *echo.Echo, appConfigUsecase domain.AppConfigUseCase, appConfigCustomValidator custom_validators.AppConfigCustomValidator, validator *validator.Validate, logger logger.ILogger) {
	// Get Localization Data
	localizationData := app_localization.ReadLocalizationFiles(consts.USER_APP)

	handler := &AppConfigHandler{
		appConfigUsecase:         appConfigUsecase,
		appConfigCustomValidator: appConfigCustomValidator,
		validator:                validator,
		logger:                   logger,
		localizationData:         localizationData,
	}
	admin := e.Group("api/v1/admin/app-config")
	{
		admin.POST("", handler.Create)
		admin.PUT("/:id", handler.Update)
		admin.GET("", handler.List)
		admin.GET("/:id", handler.FindById)
		admin.GET("/:type/by-config-type", handler.FindByType)
		admin.DELETE("/:id", handler.Delete)
	}

	mobile := e.Group("api/v1/mobile")
	{
		mobile.GET("/config", handler.FindMobileConfig)
		mobile.GET("/app-localization", handler.GetAppLocalization)
	}
	kitchen := e.Group("api/v1/kitchen")
	{
		kitchen.GET("/config", handler.FindMobileConfig)
		kitchen.GET("/app-localization", handler.GetAppLocalization)
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

	validationErr := input.Validate(ctx, a.validator, a.appConfigCustomValidator.ValidateAppTypeIsUnique())
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	var input app_config.CreateUpdateAppConfigDto
	input.Id = id
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator, a.appConfigCustomValidator.ValidateAppTypeIsUnique())
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	errResp := a.appConfigUsecase.SoftDelete(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *AppConfigHandler) FindMobileConfig(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input app_config.FindMobileConfigDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	mobileConfig, errResp := a.appConfigUsecase.FindMobileConfig(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"config": mobileConfig})
}

func (a *AppConfigHandler) GetAppLocalization(c echo.Context) error {
	return validators.SuccessResponse(c, map[string]interface{}{"localization": a.localizationData})
}
