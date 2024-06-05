package entity

import "time"

type RefreshToken struct {
	ID           int
	AccessToken  string
	RefreshToken string
	UserID       int
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	ExpiredAt    *time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
