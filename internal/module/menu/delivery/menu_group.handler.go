package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/menu/consts"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/logger"
	"samm/pkg/middlewares/portal"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type MenuGroupHandler struct {
	menuGroupUsecase domain.MenuGroupUseCase
	validator        *validator.Validate
	logger           logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitMenuGroupController(e *echo.Echo, us domain.MenuGroupUseCase, validator *validator.Validate, logger logger.ILogger, portalMiddlewares *portal.ProviderMiddlewares) {
	handler := &MenuGroupHandler{
		menuGroupUsecase: us,
		validator:        validator,
		logger:           logger,
	}
	portal := e.Group("api/v1/portal/menu-group")
	portal.Use(portalMiddlewares.AuthMiddleware)
	{
		portal.GET("", handler.List)
		portal.GET("/:id", handler.Find)
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.DELETE("/:id", handler.Delete)
		portal.PUT("/:id/change-status", handler.ChangeStatus)
		portal.DELETE("/:id/:entity/:entity_id", handler.DeleteEntity)
	}

	mobile := e.Group("api/v1/menu-group/:branch_id")
	{
		mobile.GET("/item", handler.MobileGetMenuGroupItems)
		mobile.GET("/item/:id", handler.MobileGetMenuGroupItem)
	}
}

func (a *MenuGroupHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.ListMenuGroupDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.menuGroupUsecase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"menu_groups": data})
}

func (a *MenuGroupHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.FindMenuGroupDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.menuGroupUsecase.Find(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"menu_group": data})
}

func (a *MenuGroupHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.CreateMenuGroupDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.menuGroupUsecase.Create(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *MenuGroupHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input menu_group.CreateMenuGroupDTO
	input.ID = utils.ConvertStringIdToObjectId(id)

	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id, errResp := a.menuGroupUsecase.Update(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": id})
}

func (a *MenuGroupHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.FindMenuGroupDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.menuGroupUsecase.Delete(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *MenuGroupHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input menu_group.ChangeMenuGroupStatusDto
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	if input.Id == "" {
		input.Id = id
		input.Entity = consts.MENU_CHANGE_STATUS_ENTITY
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.menuGroupUsecase.ChangeStatus(ctx, utils.ConvertStringIdToObjectId(id), &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *MenuGroupHandler) DeleteEntity(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.DeleteEntityFromMenuGroupDto
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.menuGroupUsecase.DeleteEntity(ctx, &input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *MenuGroupHandler) MobileGetMenuGroupItems(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.GetMenuGroupItemsDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.menuGroupUsecase.MobileGetMenuGroupItems(ctx, input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"menu_group_items": data})
}

func (a *MenuGroupHandler) MobileGetMenuGroupItem(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input menu_group.GetMenuGroupItemDTO
	binder := &echo.DefaultBinder{}
	if err := binder.BindHeaders(c, &input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	if err := c.Bind(&input); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	data, errResp := a.menuGroupUsecase.MobileGetMenuGroupItem(ctx, input)
	if errResp.IsError {
		return validators.ErrorResp(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"menu_group_item": data})
}
