package common

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"path/filepath"
	"samm/internal/module/common/domain"
	location "samm/internal/module/common/dto"
	"samm/pkg/config"
	"samm/pkg/logger"
	"samm/pkg/validators"
	"samm/pkg/validators/localization"
	"strconv"
	"time"
)

const tag = "CommonUseCase "

func NewCommonUseCase(repo domain.CommonRepository, logger logger.ILogger, awsS3 *s3.S3, awsConfig *config.AwsConfig) domain.CommonUseCase {
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
	awsS3     *s3.S3
	awsConfig *config.AwsConfig
}

func (l CommonUseCase) ListCities(ctx context.Context, payload *location.ListCitiesDto) (data interface{}, err validators.ErrorResponse) {

	return CitiesBuilder(payload), validators.ErrorResponse{}

}
func (l CommonUseCase) ListAssets(ctx context.Context, hasColors, hasBrands bool) (data interface{}, errResp validators.ErrorResponse) {
	assetResult := make(map[string]any)
	if hasBrands {
		carBrandsResult := ReadFile(l.logger, "/internal/module/common/consts/car_brands.json")
		assetResult["car_brands"] = carBrandsResult
	}
	if hasColors {
		carColorsResult := ReadFile(l.logger, "/internal/module/common/consts/car_colors.json")
		assetResult["car_colors"] = carColorsResult
	}
	return assetResult, validators.ErrorResponse{}

}
func (l CommonUseCase) ListCollectionMethods(ctx context.Context) (data interface{}, errResp validators.ErrorResponse) {
	collectionMethodsResult := ReadFile(l.logger, "/internal/module/common/consts/collection_methods.json")
	return collectionMethodsResult, validators.ErrorResponse{}

}
func (l CommonUseCase) FindCollectionMethodByType(ctx context.Context, collectionMethodType string) (data map[string]interface{}, errResp validators.ErrorResponse) {
	IcollectionMethodsResult, _ := l.ListCollectionMethods(ctx)
	collectionMethodsResult, ok := IcollectionMethodsResult.([]map[string]interface{})
	if !ok {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, nil)
	}
	for _, m := range collectionMethodsResult {
		if m["type"] == collectionMethodType {
			return m, validators.ErrorResponse{}
		}
	}

	return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)

}
func (l CommonUseCase) FindCollectionMethodByDefaultId(ctx context.Context, collectionMethodDefaultId string) (data map[string]interface{}, errResp validators.ErrorResponse) {
	IcollectionMethodsResult, _ := l.ListCollectionMethods(ctx)
	collectionMethodsResult, ok := IcollectionMethodsResult.([]map[string]interface{})
	if !ok {
		return nil, validators.GetErrorResponse(&ctx, localization.E1004, nil, nil)
	}
	for _, m := range collectionMethodsResult {
		if m["default_id"] == collectionMethodDefaultId {
			return m, validators.ErrorResponse{}
		}
	}

	return nil, validators.GetErrorResponse(&ctx, localization.E1002, nil, nil)

}

func (l CommonUseCase) ListCountries(ctx context.Context) (data interface{}, err validators.ErrorResponse) {

	return CountriesBuilder(), validators.ErrorResponse{}

}
func (l CommonUseCase) ReadFile(ctx context.Context, objectKey string) (string, validators.ErrorResponse) {
	//bucketName := l.awsConfig.BucketName
	//fmt.Println(objectKey, bucketName)
	//objectKey = "phase 1.jpg"
	//getObjectOutput, err := l.awsS3.GetObject(ctx, &s3.GetObjectInput{
	//	Bucket: &bucketName,
	//	Key:    &objectKey,
	//})
	//
	//if err != nil {
	//	fmt.Println(err, "GetObject", objectKey)
	//	return "", validators.GetErrorResponseFromErr(err)
	//}
	//
	//body, err := ioutil.ReadAll(getObjectOutput.Body)
	//
	//if err != nil {
	//	fmt.Println(err, "ReadAll")
	//	return "", validators.GetErrorResponseFromErr(err)
	//
	//}
	//
	//defer getObjectOutput.Body.Close()
	//return string(body), validators.ErrorResponse{}
	return string(""), validators.ErrorResponse{}
}
func (l CommonUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, filePath string) (string, validators.ErrorResponse) {

	src, err := file.Open()
	if err != nil {
		fmt.Println(err, "Open")
		return "", validators.GetErrorResponseFromErr(err)
	}
	defer src.Close()

	key := filepath.Join(filePath, strconv.Itoa(int(time.Now().Unix()))+filepath.Base(file.Filename))

	_, err = l.awsS3.PutObject(&s3.PutObjectInput{
		Body:   src,
		Bucket: aws.String(l.awsConfig.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Failed to upload data to %s/%s, %s\n", "bucket", "key", err.Error())
		return "", validators.GetErrorResponseFromErr(err)
	}

	fileLocation := filepath.Join(l.awsConfig.EndPoint, l.awsConfig.BucketName, key)
	return fileLocation, validators.ErrorResponse{}
}
