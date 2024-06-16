package event

import (
	"catalog-be/internal/domain"
	event_dto "catalog-be/internal/modules/event/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventService EventService
	validator    *validator.Validate
}

func (e *EventHandler) CreateOne(c *fiber.Ctx) error {
	body := new(event_dto.CreateEventReqeuestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(
			domain.NewErrorFiber(c, domain.NewError(400, err, nil)),
		)
	}

	if err := e.validator.Struct(body); err != nil {
		return c.Status(400).JSON(domain.NewErrorFiber(c, domain.NewError(400, err, nil)))
	}

	created, createErr := e.eventService.CreateOne(*body)
	if createErr != nil {
		return c.Status(createErr.Code).JSON(domain.NewErrorFiber(c, createErr))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": created,
	})
}

func (e *EventHandler) GetPaginatedEvents(c *fiber.Ctx) error {
	query := new(event_dto.GetEventFilter)
	if err := c.QueryParser(query); err != nil {
		return c.Status(400).JSON(
			domain.NewErrorFiber(c, domain.NewError(400, err, nil)),
		)
	}

	if err := e.validator.Struct(query); err != nil {
		return c.Status(400).JSON(domain.NewErrorFiber(c, domain.NewError(400, err, nil)))
	}

	events, getErr := e.eventService.GetPaginatedEvents(*query)
	if getErr != nil {
		return c.Status(getErr.Code).JSON(domain.NewErrorFiber(c, getErr))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     fiber.StatusOK,
		"data":     events.Data,
		"metadata": events.Metadata,
	})
}

func NewEventHandler(eventService EventService, validator *validator.Validate) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		validator:    validator,
	}
}
