package filesystem

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gbooksoft/backend/internal/database"
	"mime/multipart"
	"time"
)

func UploadMultipartFile(ctx context.Context, path string, file multipart.File) error {
	s3Client := database.GetS3Client()

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("gbook-images"),
		Key:    aws.String(path),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}

func GetPresignedURL(key string) (string, error) {
	presignClient := s3.NewPresignClient(database.GetS3Client())

	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("gbook-images"),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
