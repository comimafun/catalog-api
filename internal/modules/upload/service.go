package upload

import (
	"catalog-be/internal/domain"
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var ACCEPTED_IMAGE_TYPES = []string{
	".png",
	".jpg",
	".jpeg",
	".gif",
	".heic",
	"image/png",
	"image/jpeg",
	"image/jpeg",
	"image/heic",
	"image/gif",
}

type UploadService interface {
	randomizedFilename() (string, *domain.Error)
	validateImage(file *multipart.FileHeader) *domain.Error
	UploadImage(bucketName string, file *multipart.FileHeader) (string, *domain.Error)
}

type uploadService struct {
	s3 *s3.Client
}

// randomizedFilename implements UploadService.
func (u *uploadService) randomizedFilename() (string, *domain.Error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", domain.NewError(500, err, nil)
	}
	return uuid.String(), nil
}

// validateImage implements UploadService.
func (u *uploadService) validateImage(file *multipart.FileHeader) *domain.Error {
	var FIVE_MB = int64(5 * 1024 * 1024)
	if file.Size > FIVE_MB {
		return domain.NewError(400, errors.New("FILE_SIZE_TOO_LARGE"), nil)
	}

	if file.Header.Get("Content-Type") == "" {
		return domain.NewError(400, errors.New("FILE_TYPE_INVALID"), nil)
	}

	for _, v := range ACCEPTED_IMAGE_TYPES {
		if file.Header.Get("Content-Type") != v {
			return domain.NewError(400, errors.New("FILE_TYPE_INVALID"), nil)
		}
	}

	return nil
}

// Upload implements UploadService.
func (u *uploadService) UploadImage(bucketName string, file *multipart.FileHeader) (string, *domain.Error) {
	err := u.validateImage(file)
	if err != nil {
		return "", err
	}

	name, err := u.randomizedFilename()
	if err != nil {
		return "", domain.NewError(400, errors.New("FILE_NAME_FAILED_TO_GENERATE"), nil)
	}

	f, openErr := file.Open()
	if openErr != nil {
		return "", domain.NewError(500, openErr, nil)
	}
	defer f.Close()

	_, uploadErr := u.s3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &name,
			Body:   f,
		},
	)

	if uploadErr != nil {
		return "", domain.NewError(500, uploadErr, nil)
	}

	path := fmt.Sprintf("/%s/%s", bucketName, name)

	return path, nil
}

func NewUploadService(
	s3 *s3.Client,
) UploadService {
	return &uploadService{
		s3,
	}
}
