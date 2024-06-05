package fandom

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	fandom_dto "catalog-be/internal/modules/fandom/dto"
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FandomHandler struct {
	fandomService FandomService
	validator     *validator.Validate
}

func (h *FandomHandler) CreateOne(c *fiber.Ctx) error {
	body := new(fandom_dto.CreateBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	body.Name = strings.TrimSpace(body.Name)

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if !body.Visible {
		body.Visible = true
	}

	fandom, createErr := h.fandomService.CreateOne(entity.Fandom{Name: body.Name, Visible: body.Visible})
	if createErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorFiber(c, createErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fandom,
		"code": fiber.StatusCreated,
	})
}

func (h *FandomHandler) DeleteByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("FANDOM_ID_INVALID"), nil)))
	}

	deleteErr := h.fandomService.DeleteByID(id)
	if deleteErr != nil {
		return c.Status(deleteErr.Code).JSON(domain.NewErrorFiber(c, deleteErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": "FANDOM_DELETED",
	})
}

func (h *FandomHandler) UpdateOne(c *fiber.Ctx) error {
	body := new(fandom_dto.CreateBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	body.Name = strings.TrimSpace(body.Name)

	if err := h.validator.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, errors.New("FANDOM_ID_INVALID"), nil)))
	}

	fandom, updateErr := h.fandomService.UpdateOne(id, entity.Fandom{Name: body.Name, Visible: body.Visible})
	if updateErr != nil {
		return c.Status(updateErr.Code).JSON(domain.NewErrorFiber(c, updateErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fandom,
		"code": fiber.StatusOK,
	})
}

func (h *FandomHandler) GetFandomPagination(c *fiber.Ctx) error {
	filter := new(fandom_dto.FindAllFilter)
	if err := c.QueryParser(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if err := h.validator.Struct(filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	pagination, findErr := h.fandomService.GetPaginatedFandom(filter)
	if findErr != nil {
		return c.Status(findErr.Code).JSON(domain.NewErrorFiber(c, findErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": pagination,
		"code": fiber.StatusOK,
	})
}

func NewFandomHandler(
	fandomService FandomService,
	validator *validator.Validate,
) *FandomHandler {
	return &FandomHandler{
		fandomService,
		validator,
	}
}
