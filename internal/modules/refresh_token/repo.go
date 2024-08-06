package refreshtoken

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type RefreshTokenRepo struct {
	db *gorm.DB
}

// DeleteOneByRefreshToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteOneByRefreshToken(refreshToken string) *domain.Error {
	err := r.db.Where("token = ?", refreshToken).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteAllRefreshTokenRecordsByRefreshToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) DeleteAllRefreshTokenRecordsByRefreshToken(userID int) *domain.Error {
	err := r.db.Where("user_id = ?", userID).Unscoped().Delete(&entity.RefreshToken{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// CreateOneRefreshToken implements RefreshTokenRepo.
func (r *RefreshTokenRepo) CreateOneRefreshToken(refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
	if err := r.db.Create(&refreshToken).Error; err != nil {
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

// UpdateOneRefreshTokenByRefreshTokenID implements RefreshTokenRepo.
func (r *RefreshTokenRepo) UpdateOneRefreshTokenByRefreshTokenID(id int, refreshToken entity.RefreshToken) (*entity.RefreshToken, *domain.Error) {
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
