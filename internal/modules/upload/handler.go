package upload

import (
	"catalog-be/internal/domain"
	upload_dto "catalog-be/internal/modules/upload/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	validator     *validator.Validate
	uploadService UploadService
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	body := new(upload_dto.UploadImageBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, err, nil),
		))
	}

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, err, nil),
		))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(
			c,
			domain.NewError(fiber.StatusBadRequest, err, nil),
		))
	}

	url, uploadErr := h.uploadService.UploadImage(body.Type, file)
	if uploadErr != nil {
		return c.Status(uploadErr.Code).JSON(domain.NewErrorFiber(c, uploadErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": url,
	})
}

func NewUploadHandler(
	validator *validator.Validate,
	uploadService UploadService,
) *UploadHandler {
	return &UploadHandler{
		validator,
		uploadService,
	}
}
