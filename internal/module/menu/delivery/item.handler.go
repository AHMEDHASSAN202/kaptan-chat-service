package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"samm/internal/module/menu/custom_validators"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/item"
	"samm/pkg/logger"
	commmon "samm/pkg/middlewares/common"
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
func InitItemController(e *echo.Echo, itemUsecase domain.ItemUseCase, itemCustomValidator custom_validators.ItemCustomValidator, validator *validator.Validate, logger logger.ILogger, portalMiddlewares *portal.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &ItemHandler{
		itemUsecase:         itemUsecase,
		itemCustomValidator: itemCustomValidator,
		validator:           validator,
		logger:              logger,
	}
	portal := e.Group("api/v1/portal/item")
	portal.Use(portalMiddlewares.AuthMiddleware)
	{
		portal.POST("", handler.Create, commonMiddlewares.PermissionMiddleware("create-menus", "portal-login-accounts"))
		portal.POST("/bulk", handler.CreateBulk, commonMiddlewares.PermissionMiddleware("create-menus", "portal-login-accounts"))
		portal.GET("", handler.List, commonMiddlewares.PermissionMiddleware("list-menus", "portal-login-accounts"))
		portal.GET("/:id", handler.FindOne, commonMiddlewares.PermissionMiddleware("find-menus", "portal-login-accounts"))
		portal.PUT("/:id", handler.Update, commonMiddlewares.PermissionMiddleware("update-menus", "portal-login-accounts"))
		portal.PUT("/:id/change_status", handler.ChangeStatus, commonMiddlewares.PermissionMiddleware("update-status-menus", "portal-login-accounts"))
		portal.DELETE("/:id", handler.Delete, commonMiddlewares.PermissionMiddleware("delete-menus", "portal-login-accounts"))
	}
}

func (a *ItemHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	input := item.CreateItemDto{}
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator, a.itemCustomValidator.ValidateNameIsUnique())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.itemUsecase.Create(ctx, []item.CreateItemDto{input})
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *ItemHandler) CreateBulk(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	input := make([]item.CreateBulkItemDto, 0)
	a.logger.DumpRequest(c.Request())
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
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
	bulkInput := make([]item.CreateItemDto, 0)
	copier.Copy(&bulkInput, &input)
	errResp := a.itemUsecase.Create(ctx, bulkInput)
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
