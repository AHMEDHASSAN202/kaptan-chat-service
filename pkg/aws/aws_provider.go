package aws

import (
	"context"

	"fmt"
	pkgConfig "samm/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Init(awsConfig *pkgConfig.AwsConfig) *s3.Client {

	accessKeyID := awsConfig.AccessKey

	secretAccessKey := awsConfig.SecretKey

	cfg, err := config.LoadDefaultConfig(context.TODO(),

		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),

		config.WithRegion(awsConfig.Region),

		config.WithEndpointResolverWithOptions(

			aws.EndpointResolverWithOptionsFunc(

				func(service, region string, options ...interface{}) (aws.Endpoint, error) {

					if service == s3.ServiceID {

						return aws.Endpoint{

							URL: awsConfig.EndPoint,

							SigningRegion: awsConfig.Region,
						}, nil

					}

					return aws.Endpoint{}, &aws.EndpointNotFoundError{}

				},
			),
		),
	)

	if err != nil {

		fmt.Println("Configuration error:", err)

		return nil

	}

	client := s3.NewFromConfig(cfg)

	return client
}
