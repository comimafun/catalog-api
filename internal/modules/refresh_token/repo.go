package refreshtoken

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type RefreshTokenRepo struct {
	db *gorm.DB
}

// DeleteByRefreshToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteByRefreshToken(refreshToken string) *domain.Error {
	err := r.db.Where("token = ?", refreshToken).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteAllRecordsByUserID implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteAllRecordsByUserID(userID int) *domain.Error {
	err := r.db.Where("user_id = ?", userID).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// CreateOne implements RefreshTokenRepo.
func (r *RefreshTokenRepo) CreateOne(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	if err := r.db.Create(&refreshToken).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshToken, nil
}

// DeleteOneByAccessToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteOneByAccessToken(accessToken string) *domain.Error {
	err := r.db.Where("access_token = ?", accessToken).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteOneByID implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteOneByID(id int) *domain.Error {
	err := r.db.Delete(&entity.RefreshToken{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindOneByAccessTOken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) FindOneByAccessToken(accessToken string) (*entity.RefreshToken, *domain.Error) {
	var refreshToken entity.RefreshToken
	if err := r.db.Where("access_token = ?", accessToken).First(&refreshToken).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshToken, nil
}

// FindOneByRefreshToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) FindOneByRefreshToken(refreshToken string) (*entity.RefreshToken, *domain.Error) {
	var refreshTokenEntity entity.RefreshToken
	if err := r.db.Where("token = ?", refreshToken).First(&refreshTokenEntity).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &refreshTokenEntity, nil
}

// UpdateOneByID implements RefreshTokenRepo.
func (r *RefreshTokenRepo) UpdateOneByID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	var updated entity.RefreshToken
	if err := r.db.Model(&entity.RefreshToken{}).Where("id = ?", id).Updates(&refreshToken).Scan(&updated).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &updated, nil
}

func NewRefreshTokenRepo(db *gorm.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		db,
	}
}
