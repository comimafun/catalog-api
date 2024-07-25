package referral

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReferralHandler struct {
	service   ReferralService
	validator *validator.Validate
}

func (h *ReferralHandler) CreateReferral(c *fiber.Ctx) error {
	return nil
}

func NewReferralHandler(service ReferralService, validator *validator.Validate) *ReferralHandler {
	return &ReferralHandler{service, validator}
}
