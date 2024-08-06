package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/user/custom_validators"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/user"
	"samm/pkg/database/redis"
	echomiddleware "samm/pkg/http/echo/middleware"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
)

type UserHandler struct {
	userUsecase         domain.UserUseCase
	validator           *validator.Validate
	userCustomValidator custom_validators.UserCustomValidator
	logger              logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitUserController(e *echo.Echo, us domain.UserUseCase, validator *validator.Validate, userCustomValidator custom_validators.UserCustomValidator, logger logger.ILogger, userMiddleware *usermiddleware.Middlewares, rdb *redis.RedisClient, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &UserHandler{
		userUsecase:         us,
		validator:           validator,
		userCustomValidator: userCustomValidator,
		logger:              logger,
	}
	dashboard := e.Group("api/v1/admin/user")
	dashboard.Use(adminMiddlewares.AuthMiddleware)
	dashboard.GET("", handler.ListUser, commonMiddlewares.PermissionMiddleware("list-users"))
	dashboard.PUT("/:id/toggle-active", handler.ToggleUserActivation, commonMiddlewares.PermissionMiddleware("update-status-users"), userMiddleware.RemoveUserFromRedis)

	mobile := e.Group("api/v1/mobile/user")
	mobile.Use(echomiddleware.AppendCountryMiddleware)
	//mobile.POST("", handler.StoreUser)
	mobile.POST("/send-otp", handler.SendUserOtp)
	mobile.POST("/verify-otp", handler.VerifyUserOtp)
	mobile.POST("/sign-up", handler.SignUp, userMiddleware.AuthenticationMiddleware(true))

	authGroup := mobile.Group("", userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
	authGroup.PUT("/update-player-id", handler.UpdateUserPlayerId, userMiddleware.AuthenticationMiddleware(false))
	authGroup.PUT("", handler.UpdateUserProfile, userMiddleware.RemoveUserFromRedis)
	authGroup.GET("", handler.GetUserProfile)
	authGroup.PUT("/firebase-token-refresh", handler.RefreshFirebaseToken)
	authGroup.DELETE("", handler.DeleteUser, userMiddleware.RemoveUserFromRedis)
	authGroup.POST("/sign-out", handler.SignOut, userMiddleware.RemoveUserFromRedis)

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

	errResp := a.userUsecase.StoreUser(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *UserHandler) SendUserOtp(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.SendUserOtpDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp, otp := a.userUsecase.SendOtp(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"otp": otp})
}
func (a *UserHandler) VerifyUserOtp(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.VerifyUserOtpDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	res, errResp := a.userUsecase.VerifyOtp(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}

func (a *UserHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.UserSignUpDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	res, errResp := a.userUsecase.UserSignUp(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": res})
}

func (a *UserHandler) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.UpdateUserProfileDto
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(ctx, a.validator, a.userCustomValidator.ValidateUserEmailIsUnique())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	user, errResp := a.userUsecase.UpdateUserProfile(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, user)
}
func (a *UserHandler) GetUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	var payload dto.MobileHeaders
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	data, errResp := a.userUsecase.FindUser(&ctx, payload.CauserId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"user": data})
}

func (a *UserHandler) RefreshFirebaseToken(c echo.Context) error {
	ctx := c.Request().Context()

	var payload dto.MobileHeaders
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	data, errResp := a.userUsecase.RefreshFirebaseToken(&ctx, payload.CauserId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"firebase_token": data})
}

func (a *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	var payload dto.MobileHeaders
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	errResp := a.userUsecase.DeleteUser(&ctx, payload.CauserId)
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

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	brands, errResp := a.userUsecase.List(&ctx, &payload)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, brands)
}

func (a *UserHandler) ToggleUserActivation(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	binder := &echo.DefaultBinder{}
	var adminHeaders dto.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.userUsecase.ToggleUserActivation(&ctx, id, &adminHeaders)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *UserHandler) UpdateUserPlayerId(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.UpdateUserPlayerId
	err := c.Bind(&payload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	userDomain, errResp := a.userUsecase.UpdateUserPlayerId(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, userDomain)
}

func (a *UserHandler) SignOut(c echo.Context) error {
	ctx := c.Request().Context()

	var payload user.UserSignOutDto
	if err := c.Bind(&payload); err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	b := &echo.DefaultBinder{}
	b.BindHeaders(c, &payload)

	validationErr := payload.Validate(ctx, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.userUsecase.SignOut(&ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
