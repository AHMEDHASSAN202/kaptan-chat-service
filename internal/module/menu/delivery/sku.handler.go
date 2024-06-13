package delivery

import (
	"context"
	"samm/internal/module/menu/custom_validators"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/sku"
	"samm/pkg/logger"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SKUHandler struct {
	skuUsecase         domain.SKUUseCase
	skuCustomValidator custom_validators.SKUCustomValidator
	validator          *validator.Validate
	logger             logger.ILogger
}

// InitSKUController will initialize the article's HTTP controller
func InitSKUController(e *echo.Echo, skuUsecase domain.SKUUseCase, skuCustomValidator custom_validators.SKUCustomValidator, validator *validator.Validate, logger logger.ILogger) {
	handler := &SKUHandler{
		skuUsecase:         skuUsecase,
		skuCustomValidator: skuCustomValidator,
		validator:          validator,
		logger:             logger,
	}
	portal := e.Group("api/v1/portal/sku")
	{
		portal.POST("", handler.Create)
		portal.GET("", handler.List)
	}
}

func (a *SKUHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input sku.CreateSKUDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator, a.skuCustomValidator.ValidateSKUIsUnique())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.skuUsecase.Create(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *SKUHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input sku.ListSKUDto
	input.Pagination.SetDefault()
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	modifierGroups, errResp := a.skuUsecase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, modifierGroups)
}
