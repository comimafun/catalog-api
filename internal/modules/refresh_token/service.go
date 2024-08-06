package refreshtoken

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"catalog-be/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenService struct {
	refreshTokenRepo *RefreshTokenRepo
}

// DeleteOneByRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) DeleteOneByRefreshToken(refreshToken string) *domain.Error {
	return r.refreshTokenRepo.DeleteOneByRefreshToken(refreshToken)
}

// DeleteAllRefreshTokenRecordsByUserID implements RefreshTokenService.
func (r *RefreshTokenService) DeleteAllRefreshTokenRecordsByUserID(userID int) *domain.Error {
	return r.refreshTokenRepo.DeleteAllRefreshTokenRecordsByRefreshToken(userID)
}

// UpdateOneRefreshTokenByUserID implements RefreshTokenService.
func (r *RefreshTokenService) UpdateOneRefreshTokenByUserID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.UpdateOneRefreshTokenByRefreshTokenID(id, refreshToken)
}

// CheckSessionValidityByRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) CheckSessionValidityByRefreshToken(refreshToken string) (*entity.RefreshToken, *domain.Error) {
	token, err := r.refreshTokenRepo.FindOneByRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(401, errors.New("REFRESH_TOKEN_NOT_FOUND"), nil)
		}
		return nil, err
	}

	now := time.Now()

	if token.ExpiredAt.Before(now) {
		return nil, domain.NewError(401, errors.New("REFRESH_TOKEN_EXPIRED"), nil)
	}

	return token, nil
}

// CreateOneRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) CreateOneRefreshToken(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.CreateOneRefreshToken(refreshToken)
}

func NewRefreshTokenService(refreshTokenRepo *RefreshTokenRepo, utils utils.Utils) *RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepo,
	}
}
