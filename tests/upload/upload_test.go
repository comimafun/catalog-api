package upload_test

import (
	"bytes"
	"catalog-be/internal/database"
	"catalog-be/internal/modules/upload"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func createS3Client(ctx context.Context, t *testing.T) *s3.Client {
	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:1.4.0")
	if err != nil {
		t.Logf("Error starting localstack container: %s", err)
	}

	mappedPort, err := localstackContainer.MappedPort(ctx, "4566")
	if err != nil {
		t.Logf("Error getting mapped port: %s", err)
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		t.Logf("Error creating docker provider: %s", err)
	}
	defer provider.Close()
	host, err := provider.DaemonHost(ctx)
	if err != nil {
		t.Logf("Error getting daemon host: %s", err)
	}

	endpoint := aws.Endpoint{
		URL:               fmt.Sprintf("http://%s:%d", host, mappedPort.Int()),
		HostnameImmutable: true,
	}

	client := database.NewS3(endpoint, "us-east-1")

	return client
}

func TestMain(m *testing.M) {
	os.Setenv("BUCKET_NAME", "localstack")
	os.Setenv("APP_STAGE", "local")
	os.Setenv("ACCOUNT_KEY_ID", "test")
	os.Setenv("ACCOUNT_KEY_SECRET", "test")
	res := m.Run()
	os.Exit(res)
}

type uploadInstance struct {
	uploadService *upload.UploadService
}

func newUploadInstance(client *s3.Client) *uploadInstance {
	uploadService := upload.NewUploadService(client)
	return &uploadInstance{
		uploadService,
	}
}

func TestUpload(t *testing.T) {
	ctx := context.Background()
	client := createS3Client(ctx, t)
	// create bucket
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("localstack"),
	})
	if err != nil {
		t.Logf("Error creating bucket: %s", err)
	}

	// get bucket list
	listBucketOutput, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})

	if err != nil {
		t.Logf("Error listing buckets: %s", err)
	}

	for _, bucket := range listBucketOutput.Buckets {
		t.Logf("Bucket: %s", *bucket.Name)
	}

	instance := newUploadInstance(client)
	if client == nil {
		t.Fatal("Error creating S3 client")
	}

	t.Run("Test upload .jpeg accepted", func(t *testing.T) {
		// open the file
		file, err := os.Open("./data/accepted_jpeg.jpeg")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// create a buffer to hold the file in memory
		var buff bytes.Buffer
		buffWriter := io.Writer(&buff)

		// create a new form and create a new file field
		formWriter := multipart.NewWriter(buffWriter)
		formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		// copy the content of the file to the form's file field
		if _, err := io.Copy(formPart, file); err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
		}

		// return the multipart file header
		files, exists := multipartForm.File["file"]
		if !exists || len(files) == 0 {
			log.Fatal("multipart file not exists")
		}

		fileHeader := files[0]
		fileHeader.Header.Set("Content-Type", "image/jpeg")

		_, uploadErr := instance.uploadService.UploadImage("products", fileHeader)
		if uploadErr != nil {
			t.Logf("Error uploading image: %s", uploadErr.Err.Error())
		}
		assert.Nil(t, uploadErr)
	})
}
