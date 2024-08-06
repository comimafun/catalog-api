package report

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	report_dto "catalog-be/internal/modules/report/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	circleReportService *ReportService
	validator           *validator.Validate
}

func NewReportHandler(
	service *ReportService,
	validator *validator.Validate,
) *ReportHandler {
	return &ReportHandler{
		circleReportService: service,
		validator:           validator,
	}
}

func (rh *ReportHandler) PostCreateOneReportCircle(c *fiber.Ctx) error {
	circleID, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	var body report_dto.CreateReportPayload
	if err := c.BodyParser(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := rh.validator.Struct(&body); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	user := c.Locals("user").(*auth_dto.ATClaims)
	errInsertDb := rh.circleReportService.CreateCircleReport(circleID, user.UserID, body.Reason)
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
