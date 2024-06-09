package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/example/domain"
	"samm/internal/module/example/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type MenuHandler struct {
	menuUsecase domain.MenuUseCase
	validator   *validator.Validate
	logger      logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitUserController(e *echo.Echo, us domain.MenuUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &MenuHandler{
		menuUsecase: us,
		validator:   validator,
		logger:      logger,
	}
	g := e.Group("user")
	g.POST("", handler.Store)
	g.GET("/:id", handler.GetByID)
}

func (a *MenuHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	art, errResponse := a.menuUsecase.FindLocation(ctx, id)
	if errResponse.IsError {
		return validators.ErrorStatusBadRequest(c, errResponse)
	}

	return validators.SuccessResponse(c, art)
}

func (a *MenuHandler) Store(c echo.Context) error {
	var userRequest dto.LocationRegisterWebhook
	err := c.Bind(&userRequest)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := userRequest.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	errResp := a.menuUsecase.LocationRegisterWebhook(ctx, &dto.LocationRegisterWebhook{})
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, validationErr)
	}

	return validators.SuccessResponse(c, nil)
}
