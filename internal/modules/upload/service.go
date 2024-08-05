package upload

import (
	"catalog-be/internal/domain"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var ACCEPTED_IMAGE_TYPES = map[string]bool{
	".png":       true,
	".jpg":       true,
	".jpeg":      true,
	".gif":       true,
	".heic":      true,
	".webp":      true,
	"image/png":  true,
	"image/jpeg": true,
	"image/jpg":  true,
	"image/heic": true,
	"image/gif":  true,
	"image/webp": true,
}

var MAX_FILE_SIZE_BASED_ON_NAME = map[string]int64{
	"products":     5 * 1024 * 1024,
	"profiles":     5 * 1024 * 1024,
	"covers":       1 * 1024 * 1024,
	"descriptions": 2.5 * 1024 * 1024,
}

type UploadService struct {
	s3 *s3.Client
}

// randomizedFilename implements UploadService.
func (u *UploadService) randomizedFilename(currentName string) (string, *domain.Error) {
	uuid, err := uuid.NewV7()

	if err != nil {
		return "", domain.NewError(500, err, nil)
	}
	ext := filepath.Ext(currentName)
	uuidString := uuid.String()
	newName := fmt.Sprintf("%s%s", uuidString, ext)
	return newName, nil
}

// validateImage implements UploadService.
func (u *UploadService) validateImage(folderName string, file *multipart.FileHeader) *domain.Error {
	var MAX_FILE_SIZE = int64(5 * 1024 * 1024)
	if MAX_FILE_SIZE_BASED_ON_NAME[folderName] != 0 {
		MAX_FILE_SIZE = MAX_FILE_SIZE_BASED_ON_NAME[folderName]
	}

	if file.Size > MAX_FILE_SIZE {
		return domain.NewError(400, errors.New("FILE_SIZE_TOO_LARGE"), nil)
	}

	if file.Header.Get("Content-Type") == "" {
		return domain.NewError(400, errors.New("FILE_TYPE_INVALID"), nil)
	}

	mimeType := file.Header.Get("Content-Type")

	if !ACCEPTED_IMAGE_TYPES[mimeType] {
		return domain.NewError(400, errors.New("FILE_TYPE_INVALID"), nil)
	}

	return nil
}

// Upload implements UploadService.
func (u *UploadService) UploadImage(folderName string, file *multipart.FileHeader) (string, *domain.Error) {
	bucket := os.Getenv("BUCKET_NAME")
	appStage := os.Getenv("APP_STAGE")
	err := u.validateImage(folderName, file)
	if err != nil {
		return "", err
	}

	name, err := u.randomizedFilename(file.Filename)
	if err != nil {
		return "", domain.NewError(500, errors.New("FILE_NAME_FAILED_TO_GENERATE"), nil)
	}

	f, openErr := file.Open()
	if openErr != nil {
		return "", domain.NewError(500, openErr, nil)
	}
	defer f.Close()

	contentType := file.Header.Get("Content-Type")

	objectKey := appStage + "/" + folderName + "/" + name

	_, uploadErr := u.s3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket:      &bucket,
			Key:         &objectKey,
			Body:        f,
			ContentType: &contentType,
		},
	)

	if uploadErr != nil {
		return "", domain.NewError(500, uploadErr, nil)
	}

	path := fmt.Sprintf("/%s/%s/%s", appStage, folderName, name)

	// Final path format
	// format: /[appStage]/[foldername]/uuid.[fileExt]
	// example: /local/products/uuid.jpg
	// example: /development/profiles/uuid.jpg
	return path, nil
}

func NewUploadService(
	s3 *s3.Client,
) *UploadService {
	return &UploadService{
		s3,
	}
}
