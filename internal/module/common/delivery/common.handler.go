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
	//file, err := c.FormFile("file")
	//if err != nil {
	//	fmt.Println(err, "FormFile")
	//	return err
	//}

	return validators.SuccessResponse(c, map[string]interface{}{"data": uploadedFiles})
	//if err := c.Bind(data); err != nil {
	//	return err
	//}

	//dst, err := os.Create(file.Filename)
	//if err != nil {
	//	fmt.Println(err, "Create")
	//	return err
	//}
	//defer dst.Close()
	//uploader := manager.NewUploader(a.awsS3)
	//
	//result, err := uploader.Upload(c.Request().Context(), &s3.PutObjectInput{
	//
	//	Bucket: aws.String(a.awsConfig.BucketName),
	//
	//	Key: aws.String(filepath.Base(file.Filename)),
	//
	//	Body: src,
	//})
	//
	//if err != nil {
	//
	//	fmt.Println("Error uploading file:", err)
	//
	//	return validators.ErrorStatusBadRequest(c, validators.GetErrorResponseFromErr(err))
	//
	//}
	//
	//fmt.Printf("File uploaded successfully: %s\n", result.Location)
	//
	//return c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with usage: %s", file.Filename, data.Module))

	//ctx := c.Request().Context()
	//result, errResp := a.commonUseCase.ListCountries(ctx)
	//if errResp.IsError {
	//	a.logger.Error(errResp)
	//	return validators.ErrorStatusBadRequest(c, errResp)
	//}
	//return validators.SuccessResponse(c, map[string]interface{}{"data": result})
}

func getRelPath(location string) string {
	parse, err := url.Parse(location)
	if err != nil {
		return ""
	}
	return parse.Path
}
