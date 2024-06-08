package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/domain"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type RetailsHandler struct {
	retailsUsecase domain.LocationUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitController(e *echo.Echo, us domain.LocationUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &RetailsHandler{
		retailsUsecase: us,
		validator:      validator,
		logger:         logger,
	}
	dashboard := e.Group("admin/retails")
	dashboard.POST("", handler.StoreLocation)
}

func (a *RetailsHandler) StoreLocation(c echo.Context) error {
	a.logger.Info("Test logger")
	resp := map[string]interface{}{
		"status":  true,
		"message": "Success",
		"data":    nil,
	}
	return validators.Success(c, resp)
}
