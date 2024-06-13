package delivery

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/url"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/logger"
	"samm/pkg/validators"
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

func getRelPath(location string) string {
	parse, err := url.Parse(location)
	if err != nil {
		return ""
	}
	return parse.Path
}
