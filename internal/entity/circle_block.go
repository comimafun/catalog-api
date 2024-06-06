package entity

import (
	"time"

	"gorm.io/gorm"
)

type CircleBlock struct {
	ID        int        `json:"id"`
	Prefix    string     `json:"prefix"`
	Postfix   string     `json:"postfix"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

func (c *CircleBlock) TableName() string {
	return "circle_block"
}
