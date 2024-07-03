package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/payment/domain"
	"samm/internal/module/payment/dto/card"
	"samm/pkg/logger"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type CardHandler struct {
	cardUseCase    domain.CardUseCase
	validator      *validator.Validate
	logger         logger.ILogger
	userMiddleware *usermiddleware.Middlewares
}

// InitUserController will initialize the article's HTTP controller
func InitCardController(e *echo.Echo, us domain.CardUseCase, validator *validator.Validate, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares) {
	handler := &CardHandler{
		cardUseCase:    us,
		validator:      validator,
		logger:         logger,
		userMiddleware: userMiddleware,
	}
	mobile := e.Group("api/v1/mobile/payment/card")
	mobile.GET("", handler.ListCard, userMiddleware.AuthMiddleware)
	mobile.DELETE("/:id", handler.DeleteCard, userMiddleware.AuthMiddleware)
	mobile.GET("/:id", handler.FindCard, userMiddleware.AuthMiddleware)
}

func (a *CardHandler) FindCard(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	var payload dto.MobileHeaders
	c.Bind(&payload)

	userId := payload.CauserId
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
	var payload dto.MobileHeaders
	c.Bind(&payload)

	userId := payload.CauserId
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

	payload.UserId = payload.CauserId
	result, paginationResult, errResp := a.cardUseCase.ListCard(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result, "meta": paginationResult})
}
