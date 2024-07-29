package event

import (
	"catalog-be/internal/database/factory"
	"catalog-be/internal/domain"
	"catalog-be/internal/dto"
	"catalog-be/internal/entity"
	event_dto "catalog-be/internal/modules/event/dto"
	"catalog-be/internal/utils"
	"errors"
	"time"
)

type EventService struct {
	eventRepository *EventRepo
	utils           utils.Utils
}

// CreateOne implements EventService.
func (e *EventService) CreateOne(body event_dto.CreateEventReqeuestBody) (*entity.Event, *domain.Error) {
	startedAt, err := time.Parse(time.RFC3339, body.StartedAt)
	if err != nil {
		return nil, domain.NewError(400, errors.New("INVALID_TIME_FORMAT"), nil)
	}

	endedAt, err := time.Parse(time.RFC3339, body.EndedAt)
	if err != nil {
		return nil, domain.NewError(400, errors.New("INVALID_TIME_FORMAT"), nil)
	}

	if startedAt.After(endedAt) {
		return nil, domain.NewError(400, errors.New("INVALID_TIME_RANGE"), nil)
	}

	slug, err := e.utils.Slugify(body.Name)
	if err != nil {
		return nil, domain.NewError(500, errors.New("FAILED_TO_CREATE_SLUG"), nil)
	}

	payload := new(entity.Event)

	if body.Description == nil {
		payload.Description = ""
	}
	payload.Name = body.Name
	payload.StartedAt = startedAt
	payload.EndedAt = endedAt
	payload.Slug = slug

	created, createErr := e.eventRepository.CreateOne(*payload)
	if createErr != nil {
		return nil, domain.NewError(500, errors.New("INTERNAL_SERVER_ERROR"), nil)
	}

	return created, nil
}

// GetPaginatedEvents implements EventService.
func (e *EventService) GetPaginatedEvents(filter event_dto.GetEventFilter) (*dto.Pagination[[]entity.Event], *domain.Error) {
	count, countErr := e.eventRepository.FindAllCount(filter)
	if countErr != nil {
		return nil, domain.NewError(countErr.Code, countErr.Err, nil)
	}

	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	events, findErr := e.eventRepository.FindAll(filter)
	if findErr != nil {
		return nil, domain.NewError(findErr.Code, findErr.Err, nil)
	}

	return &dto.Pagination[[]entity.Event]{Data: events, Metadata: *metadata}, nil
}

func NewEventService(eventRepository *EventRepo, utils utils.Utils) *EventService {
	return &EventService{eventRepository: eventRepository, utils: utils}
}
