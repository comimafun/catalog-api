package upload

import (
	"catalog-be/internal/domain"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	validator     *validator.Validate
	uploadService UploadService
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	// webDomain := os.Getenv("DOMAIN")
	// if webDomain == "" {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorFiber(
	// 		c,
	// 		domain.NewError(fiber.StatusInternalServerError, errors.New("DOMAIN_IS_REQUIRED"), nil),
	// 	))
	// }

	folder := c.FormValue("type")
	errs := h.validator.Var(folder, "required,oneof=covers products profiles")
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

	// url := fmt.Sprintf("%s", path)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": path,
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
