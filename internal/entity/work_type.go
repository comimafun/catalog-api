package entity

import (
	"time"

	"gorm.io/gorm"
)

type WorkType struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (WorkType) TableName() string {
	return "work_type"
}
