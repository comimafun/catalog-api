package circle

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	circle_dto "catalog-be/internal/modules/circle/dto"
	circleblock "catalog-be/internal/modules/circle_block"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CircleHandler struct {
	circleService      CircleService
	validator          *validator.Validate
	circleBlockService circleblock.CircleBlockService
}

func (h *CircleHandler) PublishUnpublishCircle(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)

	if user.CircleID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("USER_DONT_HAVE_CIRCLE"), nil)))
	}

	publish, err := h.circleService.PublishCircleByID(*user.CircleID)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": publish,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) UpdateCircle(c *fiber.Ctx) error {

	user := c.Locals("user").(*auth_dto.ATClaims)
	if user.CircleID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("USER_DONT_HAVE_CIRCLE"), nil)))
	}

	var body circle_dto.UpdateCircleRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	circle, err := h.circleService.UpdateCircleByID(*user.CircleID, &body)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) OnboardNewCircle(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)

	var body circle_dto.OnboardNewCircleRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	body.Name = strings.TrimSpace(body.Name)

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if user.CircleID != nil {
		return c.Status(fiber.StatusConflict).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusConflict, errors.New("CIRCLE_ALREADY_EXIST"), nil)))
	}

	circle, err := h.circleService.OnboardNewCircle(&body, user.UserID)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusCreated,
	})

}

func (h *CircleHandler) FindCircleBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("SLUG_IS_EMPTY"), nil)))
	}

	circle, err := h.circleService.FindCircleBySlug(slug)
	if err != nil {
		return c.Status(err.Code).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func NewCircleHandler(
	circleService CircleService,
	validator *validator.Validate,
	circleBlockService circleblock.CircleBlockService,
) *CircleHandler {
	return &CircleHandler{
		circleService:      circleService,
		validator:          validator,
		circleBlockService: circleBlockService,
	}
}