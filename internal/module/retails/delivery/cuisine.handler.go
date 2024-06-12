package delivery

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/cuisine"
	"samm/pkg/logger"
	"samm/pkg/utils"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
)

type CuisineHandler struct {
	cuisineUsecase domain.CuisineUseCase
	validator      *validator.Validate
	logger         logger.ILogger
}

func InitCuisineController(e *echo.Echo, cuisineUsecase domain.CuisineUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &CuisineHandler{
		cuisineUsecase: cuisineUsecase,
		validator:      validator,
		logger:         logger,
	}
	portal := e.Group("api/v1/admin/cuisine")
	{
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.GET("", handler.List)
		portal.GET("/:id", handler.Find)
		portal.PUT("/:id/status", handler.ChangeStatus)
		portal.DELETE("/:id", handler.Delete)
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
	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.cuisineUsecase.Create(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}

func (a *CuisineHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input cuisine.UpdateCuisineDto
	input.Id = id

	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}

	validationErr := input.Validate(c, a.validator)
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

func (a *CuisineHandler) List(c echo.Context) error {
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
	fmt.Print("input===>", input.Query)
	input.Pagination.SetDefault()
	cuisines, paginationMeta, errResp := a.cuisineUsecase.List(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{"docs": cuisines, "meta": paginationMeta})
}

func (a *CuisineHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil))
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
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	errResp := a.cuisineUsecase.SoftDelete(&ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.SuccessResponse(c, map[string]interface{}{})
}
