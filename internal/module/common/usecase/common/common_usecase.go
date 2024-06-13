package common

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
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
func (l CommonUseCase) ListCountries(ctx context.Context) (data interface{}, err validators.ErrorResponse) {

	return CountriesBuilder(), validators.ErrorResponse{}

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
