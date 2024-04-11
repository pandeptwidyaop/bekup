package upload

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pandeptwidyaop/bekup/internal/config"
	"github.com/pandeptwidyaop/bekup/internal/exception"
	"github.com/pandeptwidyaop/bekup/internal/log"
	"github.com/pandeptwidyaop/bekup/internal/models"
)

func S3Upload(ctx context.Context, f *models.BackupFileInfo, d config.ConfigDestination) *models.BackupFileInfo {
	file, err := os.Open(f.ZipPath)
	if err != nil {
		f.Err = err
		return f
	}
	defer file.Close()

	log.GetInstance().Info("s3: uploading ", f.ZipPath, " to s3")

	client, err := newS3Client(ctx, d)
	if err != nil {
		f.Err = err
		return f
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.AWSBucket),
		Key:    aws.String(fmt.Sprintf("%s/%s/%s", d.RootDirectory, f.DatabaseName, f.ZipName)),
		Body:   file,
	})

	if err != nil {
		f.Err = err

		return f
	}

	log.GetInstance().Info("s3: ", f.ZipPath," uploaded to s3")

	return f
}

func newS3Client(ctx context.Context, src config.ConfigDestination) (*s3.Client, error) {
	if src.AWSRegion == "" {
		return nil, exception.ErrAwsRegionNotExist
	}

	if src.AWSUrl == "" {
		return nil, exception.ErrAwsUrlNotExist
	}

	if src.AWSAccessKey == "" {
		return nil, exception.ErrAwsAccessKeyNotExist
	}

	if src.AWSSecretKey == "" {
		return nil, exception.ErrAwsAccessKeySecretNotExist
	}

	if src.AWSBucket == "" {
		return nil, exception.ErrAwsBucketNotExist
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               src.AWSUrl,    // Your MinIO server's URL
				SigningRegion:     src.AWSRegion, // The region to sign requests with (use a dummy if not relevant)
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(src.AWSRegion), // Use a dummy region if not relevant
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(src.AWSAccessKey, src.AWSSecretKey, "")), // Your MinIO credentials
		awsConfig.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, err
	}

	// Create an S3 service client
	cl := s3.NewFromConfig(cfg)

	return cl, nil
}
