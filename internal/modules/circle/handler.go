package circle

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	"catalog-be/internal/modules/circle/bookmark"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"catalog-be/internal/modules/user"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CircleHandler struct {
	circleService CircleService
	validator     *validator.Validate
	userService   user.UserService
}

func (h *CircleHandler) PublishUnpublishCircle(c *fiber.Ctx) error {
	circleID, parserr := c.ParamsInt("circleid")
	if parserr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, parserr, nil)))
	}
	user := c.Locals("user").(*auth_dto.ATClaims)

	if user.CircleID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("USER_DONT_HAVE_CIRCLE"), nil)))
	}

	if *user.CircleID != circleID {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusForbidden, errors.New("FORBIDDEN"), nil)))
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
	circleID, parserr := c.ParamsInt("circleid")
	if parserr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("CIRCLE_ID_SHOULD_BE_NUMBER"), nil)))
	}
	user := c.Locals("user").(*auth_dto.ATClaims)
	if user.CircleID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("USER_DONT_HAVE_CIRCLE"), nil)))
	}

	if *user.CircleID != circleID {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusForbidden, errors.New("FORBIDDEN"), nil)))
	}

	var body circle_dto.UpdateCircleRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	circle, err := h.circleService.UpdateCircleByID(user.UserID, *user.CircleID, &body)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) OnboardNewCircle(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)

	checkUser, checkErr := h.userService.FindOneByID(user.UserID)

	if checkErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorFiber(c, checkErr))
	}

	if checkUser.CircleID != nil {
		return c.Status(fiber.StatusConflict).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusConflict, errors.New("USER_ALREADY_HAVE_CIRCLE"), nil)))
	}

	var body circle_dto.OnboardNewCircleRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	body.Name = strings.TrimSpace(body.Name)
	body.ReferralCode = strings.ToUpper(strings.TrimSpace(body.ReferralCode))

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	circle, err := h.circleService.OnboardNewCircle(&body, user.UserID)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
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

	userID := 0
	user := c.Locals("user")
	if user != nil {
		userID = user.(*auth_dto.ATClaims).UserID
	}

	circle, err := h.circleService.FindCircleBySlug(slug, userID)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) GetCircleReferral(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("circleid")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("CIRCLE_ID_SHOULD_BE_NUMBER"), nil)))
	}

	referral, refErr := h.circleService.FindReferralCodeByCircleID(circleID)

	if refErr != nil {
		if refErr.Code == 404 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"data": nil,
				"code": fiber.StatusOK,
			})
		} else {
			return c.Status(refErr.Code).JSON(domain.NewErrorFiber(c, refErr))
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": referral.ReferralCode,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) GetPaginatedCircle(c *fiber.Ctx) error {
	var query circle_dto.FindAllCircleFilter
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user")

	userID := 0
	if user != nil {
		userID = user.(*auth_dto.ATClaims).UserID
	}

	circles, err := h.circleService.GetPaginatedCircle(&query, userID)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     fiber.StatusOK,
		"data":     circles.Data,
		"metadata": circles.Metadata,
	})
}

func (h *CircleHandler) GetPaginatedBookmarkedCircle(c *fiber.Ctx) error {
	var query circle_dto.FindAllCircleFilter
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)

	circles, err := h.circleService.GetPaginatedBookmarkedCircle(user.UserID, &query)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     fiber.StatusOK,
		"data":     circles.Data,
		"metadata": circles.Metadata,
	})
}

func (h *CircleHandler) SaveCircle(c *fiber.Ctx) error {

	circleID, parseErr := c.ParamsInt("id")

	if parseErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, parseErr, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)

	err := h.circleService.SaveBookmarkCircle(circleID, user.UserID)

	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": "BOOKMARKED",
	})
}

func (h *CircleHandler) UnsaveCircle(c *fiber.Ctx) error {
	circleID, parseErr := c.ParamsInt("id")

	if parseErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, parseErr, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)

	err := h.circleService.DeleteBookmarkCircle(circleID, user.UserID)

	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": "UNBOOKMARKED",
	})
}

func (h *CircleHandler) UpdateCircleEventAttending(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)

	if user.CircleID != nil {
		if *user.CircleID != circleID {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("FORBIDDEN"), nil)))
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("FORBIDDEN"), nil)))
	}

	var body circle_dto.UpdateCircleAttendingEvent
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	circle, circleErr := h.circleService.UpdateCircleAttendingEventByID(circleID, user.UserID, &body)
	if circleErr != nil {
		return c.Status(circleErr.Code).JSON(domain.NewErrorFiber(c, circleErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func (h *CircleHandler) DeleteCircleEventAttending(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)

	if user.CircleID != nil {
		if *user.CircleID != circleID {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("FORBIDDEN"), nil)))
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("FORBIDDEN"), nil)))
	}

	circle, circleErr := h.circleService.DeleteCircleEventAttendingByID(circleID, user.UserID)
	if circleErr != nil {
		return c.Status(circleErr.Code).JSON(domain.NewErrorFiber(c, circleErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": circle,
		"code": fiber.StatusOK,
	})
}

func NewCircleHandler(
	circleService CircleService,
	validator *validator.Validate,
	userService user.UserService,
	circleBookmarkService bookmark.CircleBookmarkService,
) *CircleHandler {
	return &CircleHandler{
		circleService: circleService,
		validator:     validator,
		userService:   userService,
	}
}
