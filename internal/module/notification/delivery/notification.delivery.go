package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/notification/consts"
	"samm/internal/module/notification/domain"
	"samm/internal/module/notification/dto/notification"
	echomiddleware "samm/pkg/http/echo/middleware"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/validators"
)

type NotificationHandler struct {
	notificationUsecase domain.NotificationUseCase
	validator           *validator.Validate
	logger              logger.ILogger
}

// InitNotificationController will initialize the article's HTTP controller
func InitNotificationController(e *echo.Echo, us domain.NotificationUseCase, validator *validator.Validate, logger logger.ILogger, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares, userMiddleware *usermiddleware.Middlewares) {
	handler := &NotificationHandler{
		notificationUsecase: us,
		validator:           validator,
		logger:              logger,
	}
	dashboard := e.Group("api/v1/admin/notification")
	dashboard.Use(adminMiddlewares.AuthMiddleware)

	dashboard.POST("", handler.CreateNotification, commonMiddlewares.PermissionMiddleware("create-notifications"))
	dashboard.GET("", handler.ListNotification, commonMiddlewares.PermissionMiddleware("list-notifications"))
	dashboard.GET("/:id", handler.FindNotification, commonMiddlewares.PermissionMiddleware("find-notifications"))
	dashboard.DELETE("/:id", handler.DeleteNotification, commonMiddlewares.PermissionMiddleware("delete-notifications"))

	mobile := e.Group("api/v1/mobile/user/notification")
	mobile.Use(echomiddleware.AppendCountryMiddleware)

	mobile.GET("", handler.ListNotificationMobile, userMiddleware.AuthenticationMiddleware(false))

}
func (a *NotificationHandler) CreateNotification(c echo.Context) error {
	ctx := c.Request().Context()

	var payload notification.StoreNotificationDto

	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	if payload.Type == consts.TYPE_PRIVATE && len(payload.UserIds) == 0 {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validators.ErrorResponse{
			IsError:          true,
			ValidationErrors: map[string][]string{"user_ids": {"Required If Type Private"}},
		})
	}

	errResp := a.notificationUsecase.CreateNotification(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *NotificationHandler) FindNotification(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	data, errResp := a.notificationUsecase.FindNotification(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"notification": data})
}

func (a *NotificationHandler) DeleteNotification(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	errResp := a.notificationUsecase.DeleteNotification(ctx, id)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *NotificationHandler) ListNotification(c echo.Context) error {
	ctx := c.Request().Context()
	var payload notification.ListNotificationDto
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, errResp := a.notificationUsecase.List(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, result)
}

func (a *NotificationHandler) ListNotificationMobile(c echo.Context) error {
	ctx := c.Request().Context()
	var payload notification.ListNotificationMobileDto
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	_ = c.Bind(&payload)

	payload.Pagination.SetDefault()

	result, errResp := a.notificationUsecase.ListMobile(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, result)
}
