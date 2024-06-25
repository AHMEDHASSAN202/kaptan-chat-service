package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type UserHandler struct {
	userUsecase domain.UserUseCase
	validator   *validator.Validate
	logger      logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitUserController(e *echo.Echo, us domain.UserUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &UserHandler{
		userUsecase: us,
		validator:   validator,
		logger:      logger,
	}
	dashboard := e.Group("api/v1/admin/user")
	dashboard.POST("", handler.StoreUser)
	dashboard.GET("", handler.ListUser)
	dashboard.PUT("/:id", handler.UpdateUser)
	dashboard.GET("/:id", handler.FindUser)
	dashboard.DELETE("/:id", handler.DeleteUser)
}
func (a *UserHandler) StoreUser(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.CreateUserDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.userUsecase.StoreUser(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.UpdateUserProfileDto
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
	errResp := a.userUsecase.UpdateUser(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *UserHandler) FindUser(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.userUsecase.FindUser(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"user": data})
}

func (a *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.userUsecase.DeleteUser(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *UserHandler) ListUser(c echo.Context) error {
	ctx := c.Request().Context()
	var payload user.ListUserDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, paginationResult, errResp := a.userUsecase.ListUser(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
