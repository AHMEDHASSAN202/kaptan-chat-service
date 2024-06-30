package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/menu/custom_validators"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/logger"
	"samm/pkg/middlewares/portal"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type ItemHandler struct {
	itemUsecase         domain.ItemUseCase
	itemCustomValidator custom_validators.ItemCustomValidator
	validator           *validator.Validate
	logger              logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitItemController(e *echo.Echo, itemUsecase domain.ItemUseCase, itemCustomValidator custom_validators.ItemCustomValidator, validator *validator.Validate, logger logger.ILogger, portalMiddlewares *portal.Middlewares) {
	handler := &ItemHandler{
		itemUsecase:         itemUsecase,
		itemCustomValidator: itemCustomValidator,
		validator:           validator,
		logger:              logger,
	}
	portal := e.Group("api/v1/portal/item")
	portal.Use(portalMiddlewares.AuthMiddleware)
	{
		portal.POST("", handler.Create)
		portal.GET("", handler.List)
		portal.GET("/:id", handler.FindOne)
		portal.PUT("/:id", handler.Update)
		portal.PUT("/:id/change_status", handler.ChangeStatus)
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
		validationErr := itemDoc.Validate(ctx, a.validator, a.itemCustomValidator.ValidateNameIsUnique())
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input item.UpdateItemDto
	input.Id = id

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator, a.itemCustomValidator.ValidateNameIsUnique())
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

func (a *ItemHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input item.ListItemsDto
	binder := &echo.DefaultBinder{}
	//bind header and query params
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = binder.BindQueryParams(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	items, errResp := a.itemUsecase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, items)
}

func (a *ItemHandler) FindOne(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	item, errResp := a.itemUsecase.GetById(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, item)
}

func (a *ItemHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	errResp := a.itemUsecase.SoftDelete(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *ItemHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input item.ChangeItemStatusDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	input.Id = id
	validationErr := input.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.itemUsecase.ChangeStatus(ctx, id, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
