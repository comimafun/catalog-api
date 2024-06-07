package entity

import (
	"time"

	"gorm.io/gorm"
)

// create enum day in go consisted of `first`, `second`, `both`
type Day string

const (
	First  Day = "first"
	Second Day = "second"
	Both   Day = "both"
)

type Circle struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	PictureURL   *string        `json:"picture_url"`
	FacebookURL  *string        `json:"facebook_url"`
	InstagramURL *string        `json:"instagram_url"`
	TwitterURL   *string        `json:"twitter_url"`
	Description  *string        `json:"description"`
	Batch        *int           `json:"batch"`
	Verified     bool           `json:"verified"`
	Published    bool           `json:"published"`
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`

	CircleBlockID *int `json:"circle_block_id"`
	Day           *Day `json:"day"`
}

func (Circle) TableName() string {
	return "circle"
}
