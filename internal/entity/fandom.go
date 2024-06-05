package entity

import (
	"time"

	"gorm.io/gorm"
)

type Fandom struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Visible   bool           `json:"visible"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Fandom) TableName() string {
	return "fandom"
}
