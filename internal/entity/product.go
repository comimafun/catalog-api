package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	ImageURL  string         `json:"image_url"`
	CircleID  int            `json:"circle_id"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Product) TableName() string {
	return "product"
}
