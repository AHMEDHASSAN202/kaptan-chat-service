package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/location"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type LocationHandler struct {
	locationUsecase domain.LocationUseCase
	validator       *validator.Validate
	logger          logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitController(e *echo.Echo, us domain.LocationUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &LocationHandler{
		locationUsecase: us,
		validator:       validator,
		logger:          logger,
	}
	dashboard := e.Group("api/v1/admin/location")
	dashboard.POST("", handler.StoreLocation)
	dashboard.GET("", handler.ListLocation)
	dashboard.PUT("/:id/toggle-active", handler.ToggleLocationActive)
	dashboard.PUT("/:id", handler.UpdateLocation)
	dashboard.PUT("/:id/toggle-snooze", handler.ToggleSnooze)
	dashboard.GET("/:id", handler.FindLocation)
	dashboard.DELETE("/:id", handler.DeleteLocation)
}
func (a *LocationHandler) StoreLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var payload location.StoreLocationDto
	err := c.Bind(&payload)
	if err != nil {
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
func (a *LocationHandler) UpdateLocation(c echo.Context) error {
	ctx := c.Request().Context()

	var payload location.StoreLocationDto
	err := c.Bind(&payload)
	if err != nil {
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
	errResp := a.locationUsecase.ToggleLocationStatus(ctx, id)
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
	errResp := a.locationUsecase.DeleteLocation(ctx, id)
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil))
	}

	var input location.LocationToggleSnoozeDto
	input.Id = utils.ConvertStringIdToObjectId(id)

	err := c.Bind(&input)
	if err != nil {
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
