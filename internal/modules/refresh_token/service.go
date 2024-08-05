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

// DeleteByRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) DeleteByRefreshToken(refreshToken string) *domain.Error {
	return r.refreshTokenRepo.DeleteByRefreshToken(refreshToken)
}

// DeleteAllRecordsByUserID implements RefreshTokenService.
func (r *RefreshTokenService) DeleteAllRecordsByUserID(userID int) *domain.Error {
	return r.refreshTokenRepo.DeleteAllRecordsByUserID(userID)
}

// ForceExpiredRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) ForceExpiredRefreshToken(accessToken string) *domain.Error {
	token, err := r.refreshTokenRepo.FindOneByAccessToken(accessToken)
	if err != nil {
		return err
	}

	now := time.Now()

	_, err = r.refreshTokenRepo.UpdateOneByID(token.ID, entity.RefreshToken{
		ExpiredAt: &now,
	})

	return err
}

// FindOneByAccessToken implements RefreshTokenService.
func (r *RefreshTokenService) FindOneByAccessToken(accessToken string) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.FindOneByAccessToken(accessToken)
}

// LogoutRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) LogoutRefreshToken(id int) *domain.Error {
	return r.refreshTokenRepo.DeleteOneByID(id)
}

// UpdateByID implements RefreshTokenService.
func (r *RefreshTokenService) UpdateByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.UpdateOneByID(id, refreshToken)
}

// CheckValidityByRefreshToken implements RefreshTokenService.
func (r *RefreshTokenService) CheckValidityByRefreshToken(refreshToken string) (*entity.RefreshToken, *domain.Error) {
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

// CreateOne implements RefreshTokenService.
func (r *RefreshTokenService) CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	return r.refreshTokenRepo.CreateOne(refreshToken)
}

func NewRefreshTokenService(refreshTokenRepo *RefreshTokenRepo, utils utils.Utils) *RefreshTokenService {
	return &RefreshTokenService{
		refreshTokenRepo,
	}
}
