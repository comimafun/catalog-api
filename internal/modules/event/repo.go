package event

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	event_dto "catalog-be/internal/modules/event/dto"

	"gorm.io/gorm"
)

type EventRepo struct {
	db *gorm.DB
}

// GetEventsCount implements EventRepo.
func (e *EventRepo) GetEventsCount(filter event_dto.GetPaginatedEventsFilter) (int, *domain.Error) {

	var count int64
	err := e.db.
		Model(&entity.Event{}).
		Count(&count).Error
	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// GetPaginatedEvents implements EventRepo.
func (e *EventRepo) GetPaginatedEvents(filter event_dto.GetPaginatedEventsFilter) ([]entity.Event, *domain.Error) {
	var events []entity.Event
	err := e.db.
		Limit(filter.Limit).
		Offset(filter.Limit * (filter.Page - 1)).
		Order("started_at DESC").
		Find(&events).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return events, nil
}

// DeleteOneEventByEventID implements EventRepo.
func (e *EventRepo) DeleteOneEventByEventID(id int) *domain.Error {
	err := e.db.Delete(&entity.Event{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// CreateOneEvent implements EventRepo.
func (e *EventRepo) CreateOneEvent(event entity.Event) (*entity.Event, *domain.Error) {
	err := e.db.Create(&event).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &event, nil
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{db}
}
