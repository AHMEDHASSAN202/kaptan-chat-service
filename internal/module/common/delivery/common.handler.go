package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type CommonHandler struct {
	commonUseCase domain.CommonUseCase
	validator     *validator.Validate
	logger        logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitCommonController(e *echo.Echo, us domain.CommonUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &CommonHandler{
		commonUseCase: us,
		validator:     validator,
		logger:        logger,
	}
	dashboard := e.Group("api/v1/common")
	dashboard.GET("/countries", handler.ListCountries)
	dashboard.GET("/cities", handler.ListCities)
}
func (a *CommonHandler) ListCities(c echo.Context) error {
	ctx := c.Request().Context()
	var payload location.ListCitiesDto

	_ = c.Bind(&payload)

	result, errResp := a.commonUseCase.ListCities(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}
func (a *CommonHandler) ListCountries(c echo.Context) error {
	ctx := c.Request().Context()
	result, errResp := a.commonUseCase.ListCountries(ctx)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}
