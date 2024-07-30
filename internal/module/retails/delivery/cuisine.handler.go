package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/custom_validators"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/logger"
	"samm/pkg/middlewares/admin"
	commmon "samm/pkg/middlewares/common"
	"samm/pkg/utils"
	"samm/pkg/utils/dto"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type CuisineHandler struct {
	cuisineUsecase  domain.CuisineUseCase
	customValidator custom_validators.RetailCustomValidator
	validator       *validator.Validate
	logger          logger.ILogger
}

func InitCuisineController(e *echo.Echo, cuisineUsecase domain.CuisineUseCase, validator *validator.Validate, logger logger.ILogger, customValidator custom_validators.RetailCustomValidator, adminMiddlewares *admin.ProviderMiddlewares, commonMiddlewares *commmon.ProviderMiddlewares) {
	handler := &CuisineHandler{
		cuisineUsecase:  cuisineUsecase,
		validator:       validator,
		customValidator: customValidator,
		logger:          logger,
	}
	portal := e.Group("api/v1/admin/cuisine")
	portal.Use(adminMiddlewares.AuthMiddleware)
	{
		portal.POST("", handler.Create, commonMiddlewares.PermissionMiddleware("create-cuisines"))
		portal.PUT("/:id", handler.Update, commonMiddlewares.PermissionMiddleware("update-cuisines"))
		portal.GET("", handler.ListForDashboard, commonMiddlewares.PermissionMiddleware("list-cuisines"))
		portal.GET("/:id", handler.Find, commonMiddlewares.PermissionMiddleware("find-cuisines"))
		portal.PUT("/:id/status", handler.ChangeStatus, commonMiddlewares.PermissionMiddleware("update-status-cuisines"))
		portal.DELETE("/:id", handler.Delete, commonMiddlewares.PermissionMiddleware("delete-cuisines"))
	}
	mobile := e.Group("api/v1/mobile/cuisine")
	{
		mobile.GET("", handler.ListForMobile)
	}
}

func (a *CuisineHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input cuisine.CreateCuisineDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.customValidator.ValidateCuisineNameUnique())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	cuisine, errResp := a.cuisineUsecase.Create(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"id": cuisine.ID})

}

func (a *CuisineHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	var input cuisine.UpdateCuisineDto
	input.Id = id

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	binder := &echo.DefaultBinder{}
	if err = binder.BindHeaders(c, &input); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator, a.customValidator.ValidateCuisineNameUnique())
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.cuisineUsecase.Update(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *CuisineHandler) ListForDashboard(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input cuisine.ListCuisinesDto
	binder := &echo.DefaultBinder{}
	//bind header and query params
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = binder.BindQueryParams(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	input.Pagination.SetDefault()
	res, errResp := a.cuisineUsecase.ListCuisinesForDashboard(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}

func (a *CuisineHandler) ListForMobile(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input cuisine.ListCuisinesDto
	binder := &echo.DefaultBinder{}
	//bind header and query params
	err := binder.BindHeaders(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	err = binder.BindQueryParams(c, &input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	input.Pagination.SetDefault()
	res, errResp := a.cuisineUsecase.ListCuisinesForMobile(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, res)
}

func (a *CuisineHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	cuisine, errResp := a.cuisineUsecase.Find(&ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, cuisine)
}

func (a *CuisineHandler) ChangeStatus(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil))
	}

	var input cuisine.ChangeCuisineStatusDto
	input.Id = utils.ConvertStringIdToObjectId(id)

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.cuisineUsecase.ChangeStatus(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *CuisineHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil, nil))
	}

	binder := &echo.DefaultBinder{}
	var adminHeaders dto.AdminHeaders
	if err := binder.BindHeaders(c, &adminHeaders); err != nil {
		a.logger.Error(err)
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	errResp := a.cuisineUsecase.SoftDelete(&ctx, id, &adminHeaders)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
