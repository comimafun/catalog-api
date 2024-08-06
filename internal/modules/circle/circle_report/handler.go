package circle_report

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	circle_report_dto "catalog-be/internal/modules/circle/circle_report/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CircleReportHandler struct {
	circleReportService *CircleReportService
	validator           *validator.Validate
}

func NewCircleReportHandler(
	service *CircleReportService,
	validator *validator.Validate,
) *CircleReportHandler {
	return &CircleReportHandler{
		circleReportService: service,
		validator:           validator,
	}
}

func (crh *CircleReportHandler) CreateReportCircle(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	var body circle_report_dto.CreateCircleReportPayload
	if err := c.BodyParser(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := crh.validator.Struct(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)
	errInsertDb := crh.circleReportService.CreateCircleReport(circleID, user.UserID, body.Reason)
	if errInsertDb != nil {
		return c.
			Status(errInsertDb.Code).
			JSON(domain.NewErrorFiber(c, errInsertDb))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": true,
	})
}
