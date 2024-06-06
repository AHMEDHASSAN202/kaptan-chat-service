package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/location/domain"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type LocationHandler struct {
	locationUsecase domain.LocationUseCase
	validator       *validator.Validate
	logger          logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitController(e *echo.Echo, us domain.LocationUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &LocationHandler{
		locationUsecase: us,
		validator:       validator,
		logger:          logger,
	}
	dashboard := e.Group("admin/location")
	dashboard.POST("", handler.StoreLocation)
}

func (a *LocationHandler) StoreLocation(c echo.Context) error {
	a.logger.Info("Test logger")
	resp := map[string]interface{}{
		"status":  true,
		"message": "Success",
		"data":    nil,
	}
	return validators.Success(c, resp)
}
