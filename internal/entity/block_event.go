package entity

import (
	"time"

	"gorm.io/gorm"
)

type BlockEvent struct {
	ID        int            `json:"id"`
	EventID   int            `json:"event_id"`
	CircleID  int            `json:"circle_id"`
	Prefix    string         `json:"prefix"`
	Postfix   string         `json:"postfix"`
	Name      string         `json:"name"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (BlockEvent) TableName() string {
	return "block_event"
}
