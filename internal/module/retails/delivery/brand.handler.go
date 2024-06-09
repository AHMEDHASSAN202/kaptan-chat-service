package delivery

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"samm/internal/module/retails/domain"
	"samm/internal/module/retails/dto/brand"
	"samm/pkg/logger"
	"samm/pkg/validators"
)

type BrandHandler struct {
	brandUsecase domain.BrandUseCase
	validator    *validator.Validate
	logger       logger.ILogger
}

func InitBrandController(e *echo.Echo, brandUsecase domain.BrandUseCase, validator *validator.Validate, logger logger.ILogger) {
	handler := &BrandHandler{
		brandUsecase: brandUsecase,
		validator:    validator,
		logger:       logger,
	}
	portal := e.Group("api/v1/admin/brand")
	{
		portal.POST("", handler.Create)
		portal.PUT("/:id", handler.Update)
		portal.GET("", handler.List)
		portal.GET("/:id", handler.Find)
		//portal.PUT("/:id", handler.ChangeStatus)
		portal.DELETE("/:id", handler.Delete)
	}
}

func (a *BrandHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input brand.CreateBrandDto
	err := c.Bind(&input)
	if err != nil {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	}
	validationErr := input.Validate(c, a.validator)
	if validationErr.IsError {
		a.logger.Error(validationErr)
		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
	}

	errResp := a.brandUsecase.Create(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}

func (a *BrandHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	var input brand.UpdateBrandDto
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

	errResp := a.brandUsecase.Update(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}

func (a *BrandHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var input brand.ListBrandDto
	_ = c.Bind(&input)
	//if err != nil {
	//	return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
	//}
	input.Pagination.SetDefault()
	brands, errResp := a.brandUsecase.List(&ctx, &input)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, brands)
}

func (a *BrandHandler) Find(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	brand, errResp := a.brandUsecase.GetById(&ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, brand)
}

//func (a *BrandHandler) ChangeStatus(c echo.Context) error {
//	ctx := c.Request().Context()
//	if ctx == nil {
//		ctx = context.Background()
//	}
//
//	id := c.Param("id")
//	if id == "" {
//		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, localization.E1002, nil))
//	}
//
//	var input cuisine.ChangeCuisineStatusDto
//	input.Id = utils.ConvertStringIdToObjectId(id)
//
//	err := c.Bind(&input)
//	if err != nil {
//		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponseFromErr(err))
//	}
//
//	validationErr := input.Validate(c, a.validator)
//	if validationErr.IsError {
//		a.logger.Error(validationErr)
//		return validators.ErrorStatusUnprocessableEntity(c, validationErr)
//	}
//
//	errResp := a.brandUsecase.ChangeStatus(&ctx, &input)
//	if errResp.IsError {
//		return validators.ErrorStatusBadRequest(c, errResp)
//	}
//
//	return validators.Success(c, map[string]interface{}{})
//}

func (a *BrandHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	if id == "" {
		return validators.ErrorStatusUnprocessableEntity(c, validators.GetErrorResponse(&ctx, "E1002", nil))
	}

	errResp := a.brandUsecase.SoftDelete(&ctx, id)
	if errResp.IsError {
		return validators.ErrorStatusBadRequest(c, errResp)
	}

	return validators.Success(c, map[string]interface{}{})
}
