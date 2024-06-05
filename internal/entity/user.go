package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                int
	Name              string
	Email             string
	Hash              string
	ProfilePictureURL string
	CircleID          *int
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	DeletedAt         gorm.DeletedAt
}

func (User) TableName() string {
	return "user"
}
