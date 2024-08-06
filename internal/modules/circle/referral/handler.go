package referral

import (
	"catalog-be/internal/domain"
	referral_dto "catalog-be/internal/modules/circle/referral/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReferralHandler struct {
	service   *ReferralService
	validator *validator.Validate
}

func (h *ReferralHandler) CreateOneReferral(c *fiber.Ctx) error {
	var body referral_dto.CreateReferralPayload

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)),
		)
	}

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)),
		)
	}

	referral, err := h.service.CreateOneReferral(body)
	if err != nil {
		return c.Status(err.Code).JSON(
			domain.NewErrorFiber(c, err),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": referral,
	})
}

func NewReferralHandler(service *ReferralService, validator *validator.Validate) *ReferralHandler {
	return &ReferralHandler{service, validator}
}
