package image_optimization

import (
	"bytes"
	"catalog-be/internal/domain"
	image_optimization_dto "catalog-be/internal/modules/image_optimization/dto"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/h2non/bimg"
)

type ImageOptimizationService struct {
	s3 *s3.Client
}

func (i *ImageOptimizationService) imageExists(cacheKey string) bool {
	cdnURL := os.Getenv("CDN_URL")
	cdnURL = fmt.Sprintf("%s/%s", cdnURL, cacheKey)
	req, err := http.NewRequest("HEAD", cdnURL, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}
	var cdnClient http.Client
	resp, err := cdnClient.Do(req)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (i *ImageOptimizationService) generateKey(dto *image_optimization_dto.OptimizeImageRequest) string {
	cacheKey := "cache"
	width := "width=" + strconv.Itoa(dto.Width)
	height := "height=" + strconv.Itoa(dto.Height)
	quality := "quality=" + strconv.Itoa(dto.Quality)
	formats := []string{
		width,
		height,
		quality,
	}
	joined := strings.Join(formats[:], ",")
	cacheKey = cacheKey + "/" + joined + dto.Path
	return cacheKey
}

func (i *ImageOptimizationService) createFileHeader(path string, contentType string) (*multipart.FileHeader, *domain.Error) {
	// open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	defer file.Close()

	// create a buffer to hold the file in memory
	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	// copy the content of the file to the form's file field
	if _, err := io.Copy(formPart, file); err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	// close the form writer after the copying process is finished
	// I don't use defer in here to avoid unexpected EOF error
	formWriter.Close()

	// transform the bytes buffer into a form reader
	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, domain.NewError(500, err, nil)
	}

	fileHeader := files[0]
	fileHeader.Header.Set("Content-Type", contentType)

	return fileHeader, nil
}

func (i *ImageOptimizationService) downloadAndProcessImage(path string, dto *image_optimization_dto.OptimizeImageRequest) ([]byte, *domain.Error) {
	cdnURL := os.Getenv("CDN_URL")
	cdnURL = fmt.Sprintf("%s%s", cdnURL, path)
	fmt.Printf("cdnURL: %s\n", cdnURL)
	resp, err := http.Get(cdnURL)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewError(500, err, nil)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	newImg, newErr := i.resizeImage(imageData, dto)
	if newErr != nil {
		return nil, newErr
	}
	tempPath := filepath.Join("./temp", path)
	if err := os.MkdirAll(filepath.Dir(tempPath), os.ModePerm); err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	err = os.WriteFile(tempPath, newImg, 0644)
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return newImg, nil
}

func (i *ImageOptimizationService) resizeImage(image []byte, dto *image_optimization_dto.OptimizeImageRequest) ([]byte, *domain.Error) {
	img := bimg.NewImage(image)
	resized, err := img.Process(bimg.Options{
		Width:   dto.Width,
		Height:  dto.Height,
		Quality: dto.Quality,
	})
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return resized, nil
}

func (i *ImageOptimizationService) OptimizeImage(dto *image_optimization_dto.OptimizeImageRequest) (string, *domain.Error) {
	bucketName := os.Getenv("BUCKET_NAME")

	cacheKey := i.generateKey(dto)
	if i.imageExists(cacheKey) {
		fmt.Printf("image exists\n")
		return cacheKey, nil
	}

	data, err := i.downloadAndProcessImage(dto.Path, dto)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(data)

	_, uploadErr := i.s3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket:      &bucketName,
			Key:         &cacheKey,
			Body:        bytes.NewReader(data),
			ContentType: &contentType,
		},
	)

	if uploadErr != nil {
		return "", domain.NewError(500, uploadErr, nil)
	}

	// Defer cleanup of the temp directory
	defer func() {
		if err := os.RemoveAll("./temp"); err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}()

	fmt.Printf("image uploaded\n")
	return cacheKey, nil
}

func NewImageOptimizationService(s3 *s3.Client) *ImageOptimizationService {
	return &ImageOptimizationService{s3: s3}
}
