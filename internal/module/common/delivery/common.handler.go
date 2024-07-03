package delivery

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/url"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"strconv"
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
	dashboard.POST("/image-uploader", handler.Upload)
	dashboard.GET("/read-file", handler.ReadFile)
	dashboard.GET("/list-assets", handler.ListAssets)
	dashboard.GET("/collection-methods", handler.ListCollectionMethods)
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

func (a *CommonHandler) ReadFile(c echo.Context) error {
	ctx := c.Request().Context()
	fileName := c.QueryParam("fileName")
	result, errResp := a.commonUseCase.ReadFile(ctx, fileName)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return c.File(result)
}

func (a *CommonHandler) Upload(c echo.Context) error {
	ctx := c.Request().Context()

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println(err, "MultipartForm")
		return err
	}
	files := form.File["files"]
	uploadedFiles := []map[string]string{}
	data := location.UploadFile{}
	err = c.Bind(&data)
	if err != nil {
		fmt.Println(err, "Bind")
		return err
	}
	data.SetDefault()
	for _, file := range files {

		location, errResp := a.commonUseCase.UploadFile(ctx, file, data.Module)
		if errResp.IsError {
			a.logger.Error(errResp)
			return validators.ErrorStatusBadRequest(c, errResp)
		}
		uploadedFiles = append(uploadedFiles, map[string]string{"fullPath": location, "relPath": getRelPath(location)})
	}

	return validators.SuccessResponse(c, uploadedFiles)
}
func (a *CommonHandler) ListAssets(c echo.Context) error {
	ctx := c.Request().Context()

	IhasColors := c.QueryParam("has_colors")
	IhasBrands := c.QueryParam("has_brands")
	hasColors, err1 := strconv.ParseBool(IhasColors)
	hasBrands, err2 := strconv.ParseBool(IhasBrands)
	if err1 != nil && err2 != nil {
		a.logger.Error("hasColors or hasBrands has wrong value", err2, err1)
		return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(errors.New("hasColors or hasBrands has wrong value")))
	}

	result, errResp := a.commonUseCase.ListAssets(ctx, hasColors, hasBrands)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}

func (a *CommonHandler) ListCollectionMethods(c echo.Context) error {
	ctx := c.Request().Context()

	result, errResp := a.commonUseCase.ListCollectionMethods(ctx)
	if errResp.IsError {
		a.logger.Error(errResp)
		return validators.ErrorStatusBadRequest(c, errResp)
	}
	return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}
func getRelPath(location string) string {
	parse, err := url.Parse(location)
	if err != nil {
		return ""
	}
	return parse.Path
}
