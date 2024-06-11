package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/custom_validators"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/account"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type AccountHandler struct {
	accountUsecase        domain.AccountUseCase
	retailCustomValidator custom_validators.RetailCustomValidator
	validator             *validator.Validate
	logger                logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitAccountController(e *echo.Echo, us domain.AccountUseCase, retailCustomValidator custom_validators.RetailCustomValidator, validator *validator.Validate, logger logger.ILogger) {
	handler := &AccountHandler{
		accountUsecase:        us,
		retailCustomValidator: retailCustomValidator,
		validator:             validator,
		logger:                logger,
	}
	dashboard := e.Group("api/v1/admin/account")
	dashboard.POST("", handler.StoreAccount)
	dashboard.GET("", handler.ListAccount)
	dashboard.PUT("/:id", handler.UpdateAccount)
	dashboard.GET("/:id", handler.FindAccount)
	dashboard.DELETE("/:id", handler.DeleteAccount)
}
func (a *AccountHandler) StoreAccount(c echo.Context) error {
	ctx := c.Request().Context()

	var payload account.StoreAccountDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator, a.retailCustomValidator.ValidateAccountEmailIsUnique(""))
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.accountUsecase.StoreAccount(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) UpdateAccount(c echo.Context) error {
	ctx := c.Request().Context()

	var payload account.UpdateAccountDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	id := c.Param("id")

	validationErr := payload.Validate(c, a.validator, a.retailCustomValidator.ValidateAccountEmailIsUnique(id))
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	errResp := a.accountUsecase.UpdateAccount(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) FindAccount(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.accountUsecase.FindAccount(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"account": data})
}

func (a *AccountHandler) DeleteAccount(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.accountUsecase.DeleteAccount(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *AccountHandler) ListAccount(c echo.Context) error {
	ctx := c.Request().Context()
	var payload account.ListAccountDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, paginationResult, errResp := a.accountUsecase.ListAccount(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
