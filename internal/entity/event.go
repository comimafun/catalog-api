package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	StartedAt   time.Time      `json:"started_at"`
	EndedAt     time.Time      `json:"ended_at"`
	CreatedAt   *time.Time     `json:"-"`
	UpdatedAt   *time.Time     `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (Event) TableName() string {
	return "event"
}
