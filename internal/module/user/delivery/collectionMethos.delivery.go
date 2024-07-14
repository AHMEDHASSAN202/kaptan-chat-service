package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	common_domain "samm/internal/module/common/domain"
	"samm/internal/module/user/domain"
	"samm/internal/module/user/dto/collection_method"
	"samm/pkg/logger"
	usermiddleware "samm/pkg/middlewares/user"
	"samm/pkg/utils"
	"samm/pkg/validators"
)

type CollectionMethodHandler struct {
	collectionMethodUseCase domain.CollectionMethodUseCase
	commonUseCase           common_domain.CommonUseCase
	validator               *validator.Validate
	logger                  logger.ILogger
}

// InitUserController will initialize the article's HTTP controller
func InitCollectionMethodController(e *echo.Echo, us domain.CollectionMethodUseCase, commonUseCase common_domain.CommonUseCase, userMiddleware *usermiddleware.Middlewares, validator *validator.Validate, logger logger.ILogger) {
	handler := &CollectionMethodHandler{
		collectionMethodUseCase: us,
		commonUseCase:           commonUseCase,
		validator:               validator,
		logger:                  logger,
	}
	dashboard := e.Group("api/v1/mobile/user/collection_methods/:type")
	dashboard.POST("", handler.StoreCollectionMethod, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
	dashboard.GET("", handler.ListCollectionMethod, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
	dashboard.PUT("/:id", handler.UpdateCollectionMethod, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
	dashboard.GET("/:id", handler.FindCollectionMethod, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
	dashboard.DELETE("/:id", handler.DeleteCollectionMethod, userMiddleware.AuthenticationMiddleware(false), userMiddleware.AuthorizationMiddleware)
}

func (a *CollectionMethodHandler) StoreCollectionMethod(c echo.Context) error {
	ctx := c.Request().Context()

	var valuesPayload collection_method.Payload
	err := (&echo.DefaultBinder{}).BindBody(c, &valuesPayload)
	//err := c.Bind(&valuesPayload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	collectionMethodType := c.Param("type")
	fields, errResp := a.commonUseCase.FindCollectionMethodByType(ctx, collectionMethodType)
	if errResp.IsError {
		return validators.ErrorStatusUnprocessableEntity(c, errResp)
	}

	validationErr := collection_method.ValidatePayload(ctx, a.validator, fields, valuesPayload)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	var payload domain.CollectionMethods
	payload.Type = collectionMethodType
	payload.Fields = fields
	payload.Values = valuesPayload
	//todo: read the user_id from auth
	if v := c.Request().Header.Get(usermiddleware.CauserId); v != "" {
		payload.UserId = utils.ConvertStringIdToObjectId(v)
	}

	errResp = a.collectionMethodUseCase.StoreCollectionMethod(ctx, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *CollectionMethodHandler) UpdateCollectionMethod(c echo.Context) error {
	ctx := c.Request().Context()

	var valuesPayload collection_method.Payload
	err := (&echo.DefaultBinder{}).BindBody(c, &valuesPayload)
	//err := c.Bind(&valuesPayload)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	collectionMethodType := c.Param("type")
	id := c.Param("id")
	fields, errResp := a.commonUseCase.FindCollectionMethodByType(ctx, collectionMethodType)
	if errResp.IsError {
		return validators.ErrorStatusUnprocessableEntity(c, errResp)
	}

	validationErr := collection_method.ValidatePayload(ctx, a.validator, fields, valuesPayload)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}
	var payload domain.CollectionMethods
	payload.Type = collectionMethodType
	payload.Fields = fields
	payload.Values = valuesPayload
	//todo: read the user_id from auth
	if v := c.Request().Header.Get(usermiddleware.CauserId); v != "" {
		payload.UserId = utils.ConvertStringIdToObjectId(v)
	}

	errResp = a.collectionMethodUseCase.UpdateCollectionMethod(ctx, id, &payload)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *CollectionMethodHandler) FindCollectionMethod(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId := c.Request().Header.Get(usermiddleware.CauserId)
	//userId := "667b4dc916412390df546630"
	data, errResp := a.collectionMethodUseCase.FindCollectionMethod(ctx, id, userId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"collection_method": data})
}

func (a *CollectionMethodHandler) DeleteCollectionMethod(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId := c.Request().Header.Get(usermiddleware.CauserId)
	errResp := a.collectionMethodUseCase.DeleteCollectionMethod(ctx, id, userId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{})
}
func (a *CollectionMethodHandler) ListCollectionMethod(c echo.Context) error {
	ctx := c.Request().Context()
	collectionMethodType := c.Param("type")

	userId := c.Request().Header.Get(usermiddleware.CauserId)
	result, errResp := a.collectionMethodUseCase.ListCollectionMethod(ctx, collectionMethodType, userId)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}
