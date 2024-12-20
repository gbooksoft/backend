package filesystem

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gbooksoft/backend/internal/database"
	"mime/multipart"
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
