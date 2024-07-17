package delivery

import (
	"context"
	"github.com/jinzhu/copier"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/pkg/logger"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/middlewares/portal"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ModifierGroupHandler struct {
	modifierGroupUsecase domain.ModifierGroupUseCase
	validator            *validator.Validate
	logger               logger.ILogger
}

// InitModifierGroupController will initialize the article's HTTP controller
func InitModifierGroupController(e *echo.Echo, modifierGroupUsecase domain.ModifierGroupUseCase, validator *validator.Validate, logger logger.ILogger, portalMiddlewares *portal.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &ModifierGroupHandler{
		modifierGroupUsecase: modifierGroupUsecase,
		validator:            validator,
		logger:               logger,
	}
	portal := e.Group("api/v1/portal/modifier_group")
	portal.Use(portalMiddlewares.AuthMiddleware)
	{
		portal.POST("", handler.Create, commonMiddlewares.PermissionMiddleware("create-modifier-group", "portal-login-accounts"))
		portal.PUT("/:id", handler.Update, commonMiddlewares.PermissionMiddleware("update-modifier-group", "portal-login-accounts"))
		portal.GET("", handler.List, commonMiddlewares.PermissionMiddleware("list-modifier-group", "portal-login-accounts"))
		portal.GET("/:id", handler.Find, commonMiddlewares.PermissionMiddleware("find-modifier-group", "portal-login-accounts"))
		portal.PUT("/:id/change-status", handler.ChangeStatus, commonMiddlewares.PermissionMiddleware("update-status-modifier-group", "portal-login-accounts"))
		portal.DELETE("/:id", handler.Delete, commonMiddlewares.PermissionMiddleware("delete-modifier-group", "portal-login-accounts"))
	}
}

func (a *ModifierGroupHandler) Create(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input []modifier_group.CreateModifierGroupDto
	var inputHeaders dto.PortalHeaders
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindHeaders(c, &inputHeaders)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	for _, modifierGroupDoc := range input {
		validationErr := modifierGroupDoc.Validate(c, a.validator)
		if validationErr.IsError {
			a.logger.Error(validationErr)
			return validators.ErrorStatusUnprocessableEntity(c, validationErr)
		}
	}
	m := modifier_group.CreateUpdateModifierGroup{Data: input}
	copier.Copy(&m, &inputHeaders)
	errResp := a.modifierGroupUsecase.Create(ctx, m)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input modifier_group.UpdateModifierGroupDto
	input.Id = id

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.modifierGroupUsecase.Update(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	modifierGroup, errResp := a.modifierGroupUsecase.GetById(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, modifierGroup)
}

func (a *ModifierGroupHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input modifier_group.ListModifierGroupsDto
	input.Pagination.SetDefault()
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	modiferGroupsWithPagination, errResp := a.modifierGroupUsecase.List(ctx, &input)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"modifier_groups": modiferGroupsWithPagination})
}

func (a *ModifierGroupHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input modifier_group.ChangeModifierGroupStatusDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindPathParams(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.modifierGroupUsecase.ChangeStatus(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input modifier_group.DeleteModifierGroupDto
	err := (&echo.DefaultBinder{}).BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = (&echo.DefaultBinder{}).BindPathParams(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.modifierGroupUsecase.SoftDelete(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
