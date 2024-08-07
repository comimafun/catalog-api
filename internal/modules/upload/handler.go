package upload

import (
	"catalog-be/internal/domain"
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	validator     *validator.Validate
	uploadService *UploadService
}

func (h *UploadHandler) PostUploadImage(c *fiber.Ctx) error {
	appStage := os.Getenv("APP_STAGE")

	if appStage == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusInternalServerError, errors.New("ENV_APP_STAGE_NOT_FOUND"), nil),
		))
	}

	folder := c.FormValue("type")
	errs := h.validator.Var(folder, "required,oneof=covers products profiles descriptions")
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, errs, nil),
		))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, err, nil),
		))
	}

	if file == nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, errors.New("FILE_IS_REQUIRED"), nil),
		))
	}

	path, uploadErr := h.uploadService.UploadImage(folder, file)
	if uploadErr != nil {
		return c.Status(uploadErr.Code).JSON(domain.NewErrorFiber(c, uploadErr))
	}

	cdn := os.Getenv("CDN_URL")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": cdn + path,
	})
}

func NewUploadHandler(
	validator *validator.Validate,
	uploadService *UploadService,
) *UploadHandler {
	return &UploadHandler{
		validator,
		uploadService,
	}
}
