package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type KitchenHandler struct {
	kitchenUsecase domain.KitchenUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

// InitKitchenController will initialize the article's HTTP controller
func InitKitchenController(e *echo.Echo, us domain.KitchenUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &KitchenHandler{
		kitchenUsecase: us,
		validator:      validator,
		logger:         logger,
	}
	dashboard := e.Group("api/v1/admin/kitchen")
	dashboard.POST("", handler.CreateKitchen)
	dashboard.GET("", handler.ListKitchen)
	dashboard.PUT("/:id", handler.UpdateKitchen)
	dashboard.GET("/:id", handler.FindKitchen)
	dashboard.DELETE("/:id", handler.DeleteKitchen)
}
func (a *KitchenHandler) CreateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.StoreKitchenDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.kitchenUsecase.CreateKitchen(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) UpdateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.UpdateKitchenDto
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
	errResp := a.kitchenUsecase.UpdateKitchen(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) FindKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.kitchenUsecase.FindKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"kitchen": data})
}

func (a *KitchenHandler) DeleteKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.kitchenUsecase.DeleteKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) ListKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	var payload kitchen.ListKitchenDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, errResp := a.kitchenUsecase.List(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, result)
}

	