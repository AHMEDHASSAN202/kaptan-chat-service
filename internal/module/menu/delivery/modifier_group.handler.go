package delivery

import (
	"context"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/modifier_group"
	"samm/pkg/logger"
	"samm/pkg/validators"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ModifierGroupHandler struct {
	modifierGroupUsecase domain.ModifierGroupUseCase
	validator            *validator.Validate
	logger               logger.ILogger
}

// InitModifierGroupController will initialize the article's HTTP controller
func InitModifierGroupController(e *echo.Echo, modifierGroupUsecase domain.ModifierGroupUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &ModifierGroupHandler{
		modifierGroupUsecase: modifierGroupUsecase,
		validator:            validator,
		logger:               logger,
	}
	portal := e.Group("api/v1/portal/modifier_group")
	{
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.GET("", handler.List)
		portal.GET("/:id", handler.Find)
		portal.PUT("/:id/change-status", handler.ChangeStatus)
		portal.DELETE("/:id", handler.Delete)
	}
}

func (a *ModifierGroupHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input modifier_group.CreateUpdateModifierGroupDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.modifierGroupUsecase.Create(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input modifier_group.CreateUpdateModifierGroupDto
	input.Id = id

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	// validationErr := input.Validate(c, a.validator)
	// if validationErr.IsError {
	// 	a.logger.Error(validationErr)
	// 	return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	// }

	errResp := a.modifierGroupUsecase.Update(ctx, input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	modifierGroup, errResp := a.modifierGroupUsecase.GetById(ctx, id)
	if errResp.IsError {
		logger.Logger.Info("==================>")
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, modifierGroup)
}

func (a *ModifierGroupHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input modifier_group.ListModifierGroupsDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	modifierGroups, errResp := a.modifierGroupUsecase.List(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, modifierGroups)
}

func (a *ModifierGroupHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input modifier_group.ChangeModifierGroupStatusDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.modifierGroupUsecase.ChangeStatus(ctx, id, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}

func (a *ModifierGroupHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	errResp := a.modifierGroupUsecase.SoftDelete(ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}
