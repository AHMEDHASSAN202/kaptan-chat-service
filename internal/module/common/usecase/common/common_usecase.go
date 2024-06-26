package common

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/config"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"strconv"
	"time"
)

const tag = "CommonUseCase "

func NewCommonUseCase(repo domain.CommonRepository, logger logger.ILogger, awsS3 *s3.Client,
	awsConfig *config.AwsConfig) domain.CommonUseCase {
	return &CommonUseCase{
		repo:      repo,
		logger:    logger,
		awsS3:     awsS3,
		awsConfig: awsConfig,
	}
}

type CommonUseCase struct {
	repo      domain.CommonRepository
	logger    logger.ILogger
	awsS3     *s3.Client
	awsConfig *config.AwsConfig
}

func (l CommonUseCase) ListCities(ctx context.Context, payload *location.ListCitiesDto) (data interface{}, err validators.ErrorResponse) {

	return CitiesBuilder(payload), validators.ErrorResponse{}

}
func (l CommonUseCase) ListAssets(ctx context.Context, hasColors, hasBrands bool) (data interface{}, errResp validators.ErrorResponse) {
	assetResult := make(map[string]any)
	if hasBrands {
		pwd, _ := os.Getwd()
		carBrandsFile, err := os.Open(pwd + "/internal/module/common/consts/car_brands.json")
		if err != nil {
			l.logger.Error(err)
		}

		defer carBrandsFile.Close()
		reader := bufio.NewReader(carBrandsFile)
		content, err := io.ReadAll(reader)
		if err != nil {
			l.logger.Error(err)
		}

		var carBrandsResult []map[string]interface{}
		json.Unmarshal(content, &carBrandsResult)
		assetResult["carBrands"] = carBrandsResult
	}
	if hasColors {
		pwd, _ := os.Getwd()
		carColorsFile, err := os.Open(pwd + "/internal/module/common/consts/car_colors.json")
		if err != nil {
			l.logger.Error(err)
		}

		defer carColorsFile.Close()
		reader := bufio.NewReader(carColorsFile)
		content, err := io.ReadAll(reader)
		if err != nil {
			l.logger.Error(err)
		}
		var carColorsResult []map[string]interface{}
		json.Unmarshal(content, &carColorsResult)
		assetResult["carColors"] = carColorsResult
	}
	return assetResult, validators.ErrorResponse{}

}
func (l CommonUseCase) ListCountries(ctx context.Context) (data interface{}, err validators.ErrorResponse) {

	return CountriesBuilder(), validators.ErrorResponse{}

}
func (l CommonUseCase) ReadFile(ctx context.Context, objectKey string) (string, validators.ErrorResponse) {
	bucketName := l.awsConfig.BucketName
	fmt.Println(objectKey, bucketName)
	objectKey = "phase 1.jpg"
	getObjectOutput, err := l.awsS3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		fmt.Println(err, "GetObject", objectKey)
		return "", validators.GetErrorResponseFromErr(err)
	}

	body, err := ioutil.ReadAll(getObjectOutput.Body)

	if err != nil {
		fmt.Println(err, "ReadAll")
		return "", validators.GetErrorResponseFromErr(err)

	}

	defer getObjectOutput.Body.Close()
	return string(body), validators.ErrorResponse{}
}
func (l CommonUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, filePath string) (string, validators.ErrorResponse) {

	src, err := file.Open()
	if err != nil {
		fmt.Println(err, "Open")
		return "", validators.GetErrorResponseFromErr(err)
	}
	defer src.Close()

	uploader := manager.NewUploader(l.awsS3)

	uploaderResp, err := uploader.Upload(ctx, &s3.PutObjectInput{

		Bucket: aws.String(l.awsConfig.BucketName),

		Key: aws.String(filepath.Join(filePath, strconv.Itoa(int(time.Now().Unix()))+filepath.Base(file.Filename))),

		Body: src,
	})

	if err != nil {
		fmt.Println("Error uploading file:", err)
		return "", validators.GetErrorResponseFromErr(err)

	}
	return uploaderResp.Location, validators.ErrorResponse{}
}
