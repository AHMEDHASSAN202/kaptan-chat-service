package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/card"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type CardHandler struct {
	cardUseCase domain.CardUseCase
	validator   *validator.Validate
	logger      logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitCardController(e *echo.Echo, us domain.CardUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &CardHandler{
		cardUseCase: us,
		validator:   validator,
		logger:      logger,
	}
	mobile := e.Group("api/v1/mobile/payment/card")
	mobile.POST("", handler.CreateCard)
	mobile.GET("", handler.ListCard)
	mobile.DELETE("/:id", handler.DeleteCard)
	mobile.GET("/:id", handler.FindCard)
}

func (a *CardHandler) CreateCard(c echo.Context) error {
	ctx := c.Request().Context()

	var payload card.CreateCardDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	// TODO GET FROM AUTH
	payload.UserId = "667bb56c08d41655bda52a83"
	errResp := a.cardUseCase.StoreCard(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *CardHandler) FindCard(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId := "667bb56c08d41655bda52a83"
	data, errResp := a.cardUseCase.FindCard(ctx, id, userId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"card": data})
}

func (a *CardHandler) DeleteCard(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId := "667bb56c08d41655bda52a83"
	errResp := a.cardUseCase.DeleteCard(ctx, id, userId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *CardHandler) ListCard(c echo.Context) error {
	ctx := c.Request().Context()
	var payload card.ListCardDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	payload.UserId = "667bb56c08d41655bda52a83"
	result, paginationResult, errResp := a.cardUseCase.ListCard(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
