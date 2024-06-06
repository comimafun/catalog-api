package work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"github.com/gofiber/fiber/v2"
)

type WorkTypeHandler struct {
	workTypeService WorkTypeService
}

type createUpdateRequestBody struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

func (h *WorkTypeHandler) CreateOne(c *fiber.Ctx) error {
	var body createUpdateRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	workType, err := h.workTypeService.CreateOne(entity.WorkType{Name: body.Name})
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": workType,
		"code": fiber.StatusCreated,
	})
}

func (h *WorkTypeHandler) UpdateOne(c *fiber.Ctx) error {
	var body createUpdateRequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, paramsErr := c.ParamsInt("id")
	if paramsErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, paramsErr, nil)))
	}

	workType, err := h.workTypeService.UpdateOne(id, entity.WorkType{Name: body.Name})
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": workType,
		"code": fiber.StatusOK,
	})
}

func (h *WorkTypeHandler) DeleteByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	deleteErr := h.workTypeService.DeleteByID(id)
	if deleteErr != nil {
		return c.Status(deleteErr.Code).JSON(domain.NewErrorFiber(c, deleteErr))
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"code": fiber.StatusNoContent,
		"data": "DELETED",
	})
}

func (h *WorkTypeHandler) GetAll(c *fiber.Ctx) error {
	workTypes, err := h.workTypeService.FindAll()
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": workTypes,
	})
}

func NewWorkTypeHandler(
	workTypeService WorkTypeService,
) *WorkTypeHandler {
	return &WorkTypeHandler{
		workTypeService: workTypeService,
	}
}
