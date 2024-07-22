package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	echomiddleware "samm/pkg/http/echo/middleware"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type LocationHandler struct {
	locationUsecase domain.LocationUseCase
	validator       *validator.Validate
	logger          logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitController(e *echo.Echo, us domain.LocationUseCase, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &LocationHandler{
		locationUsecase: us,
		validator:       validator,
		logger:          logger,
	}
	dashboard := e.Group("api/v1/admin/location")
	dashboard.Use(adminMiddlewares.AuthMiddleware)

	dashboard.POST("", handler.StoreLocation, commonMiddlewares.PermissionMiddleware("create-locations-accounts"))
	dashboard.POST("/bulk", handler.BulkStoreLocation, commonMiddlewares.PermissionMiddleware("create-locations-accounts"))
	dashboard.GET("", handler.ListLocation, commonMiddlewares.PermissionMiddleware("list-locations-accounts"))
	dashboard.PUT("/:id/toggle-active", handler.ToggleLocationActive, commonMiddlewares.PermissionMiddleware("update-status-locations-accounts"))
	dashboard.PUT("/:id", handler.UpdateLocation, commonMiddlewares.PermissionMiddleware("update-locations-accounts"))
	dashboard.PUT("/:id/toggle-snooze", handler.ToggleSnooze, commonMiddlewares.PermissionMiddleware("update-status-locations-accounts"))
	dashboard.GET("/:id", handler.FindLocation, commonMiddlewares.PermissionMiddleware("find-locations-accounts"))
	dashboard.DELETE("/:id", handler.DeleteLocation, commonMiddlewares.PermissionMiddleware("delete-locations-accounts"))

	mobile := e.Group("api/v1/mobile/location")
	mobile.Use(echomiddleware.AppendCountryMiddleware)
	mobile.GET("", handler.ListMobileLocation)
	mobile.GET("/search", handler.SearchMobileLocation)
	mobile.GET("/:id", handler.FindMobileLocation)

}
func (a *LocationHandler) StoreLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var payload location.StoreLocationDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.locationUsecase.StoreLocation(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *LocationHandler) BulkStoreLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var payload location.StoreBulkLocationDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.locationUsecase.BulkStoreLocation(ctx, payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *LocationHandler) UpdateLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var payload location.StoreLocationDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &payload); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	id := c.Param("id")
	errResp := a.locationUsecase.UpdateLocation(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *LocationHandler) ToggleLocationActive(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	binder := &echo.DefaultBinder{}
	var adminHeaders dto.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.locationUsecase.ToggleLocationStatus(ctx, id, &adminHeaders)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *LocationHandler) FindLocation(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.locationUsecase.FindLocation(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"location": data})
}
func (a *LocationHandler) DeleteLocation(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	binder := &echo.DefaultBinder{}
	var adminHeaders dto.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.locationUsecase.DeleteLocation(ctx, id, &adminHeaders)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *LocationHandler) ListLocation(c echo.Context) error {
	ctx := c.Request().Context()
	var payload location.ListLocationDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, paginationResult, errResp := a.locationUsecase.ListLocation(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
func (a *LocationHandler) ToggleSnooze(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input location.LocationToggleSnoozeDto
	input.Id = utils.ConvertStringIdToObjectId(id)

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.locationUsecase.ToggleSnooze(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *LocationHandler) ListMobileLocation(c echo.Context) error {
	ctx := c.Request().Context()
	var payload location.ListLocationMobileDto
	_ = c.Bind(&payload)
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	payload.SetDefault()

	result, paginationResult, errResp := a.locationUsecase.ListMobileLocation(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
func (a *LocationHandler) SearchMobileLocation(c echo.Context) error {
	ctx := c.Request().Context()
	var payload location.ListLocationMobileDto
	_ = c.Bind(&payload)
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	payload.SetDefault()

	result, paginationResult, errResp := a.locationUsecase.ListMobileLocation(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
func (a *LocationHandler) FindMobileLocation(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var payload location.FindLocationMobileDto
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)
	b.BindQueryParams(c, &payload)
	data, errResp := a.locationUsecase.FindMobileLocation(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"location": data})
}
