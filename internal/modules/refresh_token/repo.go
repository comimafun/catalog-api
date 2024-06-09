package refreshtoken

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type RefreshTokenRepo interface {
	CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	DeleteOneByID(id int) *domain.Error
	DeleteOneByAccessToken(accessToken string) *domain.Error
	UpdateOneByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error)
	FindOneByRefreshToken(refreshToken string) (*entity.RefreshToken, *domain.Error)
	FindOneByAccessToken(accessToken string) (*entity.RefreshToken, *domain.Error)
	DeleteAllRecordsByUserID(userID int) *domain.Error
	DeleteByRefreshToken(refreshToken string) *domain.Error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

// DeleteByRefreshToken implements RefreshTokenRepo.
func (r *refreshTokenRepo) DeleteByRefreshToken(refreshToken string) *domain.Error {
	err := r.db.Where("token = ?", refreshToken).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteAllRecordsByUserID implements RefreshTokenRepo.
func (r *refreshTokenRepo) DeleteAllRecordsByUserID(userID int) *domain.Error {
	err := r.db.Where("user_id = ?", userID).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// CreateOne implements RefreshTokenRepo.
func (r *refreshTokenRepo) CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	if err := r.db.Create(&refreshToken).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshToken, nil
}

// DeleteOneByAccessToken implements RefreshTokenRepo.
func (r *refreshTokenRepo) DeleteOneByAccessToken(accessToken string) *domain.Error {
	err := r.db.Where("access_token = ?", accessToken).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteOneByID implements RefreshTokenRepo.
func (r *refreshTokenRepo) DeleteOneByID(id int) *domain.Error {
	err := r.db.Delete(&entity.RefreshToken{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindOneByAccessTOken implements RefreshTokenRepo.
func (r *refreshTokenRepo) FindOneByAccessToken(accessToken string) (*entity.RefreshToken, *domain.Error) {
	var refreshToken entity.RefreshToken
	if err := r.db.Where("access_token = ?", accessToken).First(&refreshToken).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshToken, nil
}

// FindOneByRefreshToken implements RefreshTokenRepo.
func (r *refreshTokenRepo) FindOneByRefreshToken(refreshToken string) (*entity.RefreshToken, *domain.Error) {
	var refreshTokenEntity entity.RefreshToken
	if err := r.db.Where("token = ?", refreshToken).First(&refreshTokenEntity).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshTokenEntity, nil
}

// UpdateOneByID implements RefreshTokenRepo.
func (r *refreshTokenRepo) UpdateOneByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	var updated entity.RefreshToken
	if err := r.db.Model(&entity.RefreshToken{}).Where("id = ?", id).Updates(&refreshToken).Scan(&updated).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &updated, nil
}

func NewRefreshTokenRepo(db *gorm.DB) RefreshTokenRepo {
	return &refreshTokenRepo{
		db,
	}
}
