package auth_dto

import (
	"catalog-be/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleUserData struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Name           string `json:"name"`
	Given_name     string `json:"given_name"`
	Family_name    string `json:"family_name"`
	Picture        string `json:"picture"`
	Locale         string `json:"locale"`
}

type BasicClaims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	CircleID *int   `json:"circle_id"`
}

type ATClaims struct {
	BasicClaims
	jwt.RegisteredClaims
}

type NewToken struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiredAt  time.Time
	RefreshTokenExpiredAt time.Time
}

type NewTokenResponse struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	AccessTokenExpiredAt  string `json:"access_token_expired_at"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at"`
}

type SelfResponse struct {
	User                  entity.User    `json:"user"`
	Circle                *entity.Circle `json:"circle"`
	AccessTokenExpiredAt  string         `json:"access_token_expired_at"`
	RefreshTokenExpiredAt string         `json:"refresh_token_expired_at"`
}
