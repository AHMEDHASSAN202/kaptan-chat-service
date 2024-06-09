package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type ItemHandler struct {
	itemUsecase domain.ItemUseCase
	validator   *validator.Validate
	logger      logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitItemController(e *echo.Echo, itemUsecase domain.ItemUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &ItemHandler{
		itemUsecase: itemUsecase,
		validator:   validator,
		logger:      logger,
	}
	portal := e.Group("api/v1/portal/cuisine")
	{
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.DELETE("/:id", handler.Delete)
	}
}

func (a *ItemHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input []item.CreateItemDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	for _, itemDoc := range input {
		validationErr := itemDoc.Validate(c, a.validator)
		if validationErr.IsError {
			a.logger.Error(validationErr)
			return validators.ErrorStatusUnprocessableEntity(c, validationErr)
		}
	}

	errResp := a.itemUsecase.Create(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *ItemHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input item.UpdateItemDto
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

	errResp := a.itemUsecase.Update(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *ItemHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	errResp := a.itemUsecase.SoftDelete(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
