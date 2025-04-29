package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	pkgConfig "kaptan/pkg/config"
)

func Init(awsConfig *pkgConfig.Config) *s3.S3 {
	accessKeyID := awsConfig.AwsConfig.AccessKey

	secretAccessKey := awsConfig.AwsConfig.SecretKey
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(awsConfig.AwsConfig.EndPoint),
		Region:           aws.String(awsConfig.AwsConfig.Region),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	//**just for test to read the buckets
	//result, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	//if err != nil {
	//	fmt.Printf("unable to list buckets, %v", err)
	//}
	//
	//for _, bucket := range result.Buckets {
	//	fmt.Printf("* %s\n", aws.StringValue(bucket.Name))
	//}
	//
	//fmt.Printf("Successfully uploaded %q to %q\n", "keyName", "bucketName")
	return s3Client
}
