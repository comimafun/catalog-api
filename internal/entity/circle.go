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

	Day *Day `json:"day"`
}

type CircleRaw struct {
	Circle

	FandomID        int            `json:"fandom_id"`
	FandomName      string         `json:"fandom_name"`
	FandomVisible   bool           `json:"fandom_visible"`
	FandomCreatedAt *time.Time     `json:"fandom_created_at"`
	FandomUpdatedAt *time.Time     `json:"fandom_updated_at"`
	FandomDeletedAt gorm.DeletedAt `json:"-"`
}

func (Circle) TableName() string {
	return "circle"
}

type CircleFandom struct {
	CircleID  int        `json:"circle_id"`
	FandomID  int        `json:"fandom_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (CircleFandom) TableName() string {
	return "circle_fandom"
}
