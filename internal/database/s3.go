package database

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"
)

var s3Client *s3.Client

func SetUpS3Connection(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("S3_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		)),
		config.WithBaseEndpoint(os.Getenv("S3_BASE_ENDPOINT")), // Эндпоинт Yandex S3
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	err = checkS3Connection(ctx, client)

	if err != nil {
		return nil, err
	}

	s3Client = client

	return s3Client, nil
}

func checkS3Connection(ctx context.Context, client *s3.Client) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	_, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	return err
}

func GetS3Client() *s3.Client {
	return s3Client
}
