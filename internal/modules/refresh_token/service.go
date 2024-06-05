package refreshtoken

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"catalog-be/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenService interface {
	CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	LogoutRefreshToken(id int) *domain.Error
	UpdateByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	CheckValidityByRefreshToken(refreshToken string) (bool, *domain.Error)
}

type refreshTokenService struct {
	refreshTokenRepo RefreshTokenRepo
}

// LogoutRefreshToken implements RefreshTokenService.
func (r *refreshTokenService) LogoutRefreshToken(id int) *domain.Error {
	return r.refreshTokenRepo.DeleteOneByID(id)
}

// UpdateByID implements RefreshTokenService.
func (r *refreshTokenService) UpdateByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.UpdateOneByID(id, refreshToken)
}

// CheckValidityByRefreshToken implements RefreshTokenService.
func (r *refreshTokenService) CheckValidityByRefreshToken(refreshToken string) (bool, *domain.Error) {
	token, err := r.refreshTokenRepo.FindOneByRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return false, domain.NewError(401, errors.New("REFRESH_TOKEN_NOT_FOUND"), nil)
		}
		return false, err
	}

	now := time.Now()

	if token.ExpiredAt.Before(now) {
		return false, domain.NewError(401, errors.New("REFRESH_TOKEN_EXPIRED"), nil)
	}

	return true, nil
}

// CreateOne implements RefreshTokenService.
func (r *refreshTokenService) CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.CreateOne(refreshToken)
}

func NewRefreshTokenService(refreshTokenRepo RefreshTokenRepo, utils utils.Utils) RefreshTokenService {
	return &refreshTokenService{
		refreshTokenRepo,
	}
}
