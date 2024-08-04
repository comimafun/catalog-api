package image_optimization

import (
	"catalog-be/internal/domain"
	image_optimization_dto "catalog-be/internal/modules/image_optimization/dto"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ImageOptimizationHandler struct {
	service   *ImageOptimizationService
	validator *validator.Validate
}

func (h *ImageOptimizationHandler) GetOptimizedImage(c *fiber.Ctx) error {

	var req image_optimization_dto.OptimizeImageRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	if req.Quality == 0 {
		req.Quality = 80
	}

	cachedPath, err := h.service.OptimizeImage(&req)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	cdnURL := os.Getenv("CDN_URL")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": cdnURL + "/" + cachedPath,
	})
}

func NewImageOptimizationHandler(service *ImageOptimizationService, validator *validator.Validate) *ImageOptimizationHandler {
	return &ImageOptimizationHandler{service: service, validator: validator}
}
