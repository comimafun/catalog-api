package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Hash              string         `json:"-"`
	ProfilePictureURL string         `json:"profile_picture_url"`
	CircleID          *int           `json:"circle_id"`
	CreatedAt         *time.Time     `json:"-"`
	UpdatedAt         *time.Time     `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"`
}

func (User) TableName() string {
	return "user"
}
