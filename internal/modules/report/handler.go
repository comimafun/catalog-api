package report

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	auth_dto "catalog-be/internal/modules/auth/dto"
	report_dto "catalog-be/internal/modules/report/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	reportService *ReportService
	validator     *validator.Validate
}

func NewReportHandler(
	service *ReportService,
	validator *validator.Validate,
) *ReportHandler {
	return &ReportHandler{
		reportService: service,
		validator:     validator,
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

	serviceErr := rh.reportService.CreateReportCircle(&entity.Report{
		UserID:   user.UserID,
		CircleID: circleID,
		Reason:   body.Reason,
	})
	if serviceErr != nil {
		return c.
			Status(serviceErr.Code).
			JSON(domain.NewErrorFiber(c, serviceErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": true,
	})
}
