package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	custom_validators2 "samm/internal/module/admin/custom_validators"
	"samm/internal/module/kitchen/custom_validators"
	"samm/internal/module/kitchen/domain"
	"samm/internal/module/kitchen/dto/kitchen"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/validators"
)

type KitchenHandler struct {
	kitchenUsecase         domain.KitchenUseCase
	adminCustomValidator   custom_validators2.AdminCustomValidator
	kitchenCustomValidator custom_validators.KitchenCustomValidator
	validator              *validator.Validate
	logger                 logger.ILogger
}

// InitKitchenController will initialize the article's HTTP controller
func InitKitchenController(e *echo.Echo, us domain.KitchenUseCase, validator *validator.Validate, logger logger.ILogger, adminCustomValidator custom_validators2.AdminCustomValidator, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares, kitchenCustomValidator custom_validators.KitchenCustomValidator) {
	handler := &KitchenHandler{
		kitchenUsecase:         us,
		validator:              validator,
		logger:                 logger,
		adminCustomValidator:   adminCustomValidator,
		kitchenCustomValidator: kitchenCustomValidator,
	}
	dashboard := e.Group("api/v1/admin/kitchen")
	dashboard.Use(adminMiddlewares.AuthMiddleware)

	dashboard.POST("", handler.CreateKitchen, commonMiddlewares.PermissionMiddleware("create-kitchens"))
	dashboard.GET("", handler.ListKitchen, commonMiddlewares.PermissionMiddleware("list-kitchens"))
	dashboard.PUT("/:id", handler.UpdateKitchen, commonMiddlewares.PermissionMiddleware("update-kitchens"))
	dashboard.GET("/:id", handler.FindKitchen, commonMiddlewares.PermissionMiddleware("find-kitchens"))
	dashboard.DELETE("/:id", handler.DeleteKitchen, commonMiddlewares.PermissionMiddleware("delete-kitchens"))

	//mobile_kitchen := e.Group("api/v1/mobile-kitchen")
	//mobile_kitchen.Use(adminMiddlewares.AuthMiddleware)

}
func (a *KitchenHandler) CreateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.StoreKitchenDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(c, a.validator, a.adminCustomValidator.ValidateEmailIsUnique(), a.kitchenCustomValidator.ValidateAccountAndLocationRequired())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.kitchenUsecase.CreateKitchen(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) UpdateKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	var payload kitchen.UpdateKitchenDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator, a.kitchenCustomValidator.ValidateAccountAndLocationRequired())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	id := c.Param("id")
	errResp := a.kitchenUsecase.UpdateKitchen(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) FindKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.kitchenUsecase.FindKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"kitchen": data})
}

func (a *KitchenHandler) DeleteKitchen(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.kitchenUsecase.DeleteKitchen(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *KitchenHandler) ListKitchen(c echo.Context) error {
	ctx := c.Request().Context()
	var payload kitchen.ListKitchenDto

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, errResp := a.kitchenUsecase.List(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, result)
}
