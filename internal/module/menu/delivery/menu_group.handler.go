package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/menu/domain"
	"samm/internal/module/menu/dto/menu_group"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type MenuGroupHandler struct {
	menuGroupUsecase domain.MenuGroupUseCase
	validator        *validator.Validate
	logger           logger.ILogger
}

// InitMenuGroupController will initialize the article's HTTP controller
func InitMenuGroupController(e *echo.Echo, us domain.MenuGroupUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &MenuGroupHandler{
		menuGroupUsecase: us,
		validator:        validator,
		logger:           logger,
	}
	portal := e.Group("api/v1/portal/")
	{
		portal.POST("menu-group", handler.Create)
	}
}

func (a *MenuGroupHandler) Create(c echo.Context) error {
	var input menu_group.CreateMenuGroupDTO
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id, errResp := a.menuGroupUsecase.Create(ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{"id": id})
}
