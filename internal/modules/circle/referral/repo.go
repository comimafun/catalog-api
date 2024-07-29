package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ReferralRepo struct {
	db *gorm.DB
}

// FindOneReferralByCircleID implements ReferralRepo.
func (r *ReferralRepo) FindOneReferralByCircleID(circleID int) (*entity.Referral, *domain.Error) {
	referral := &entity.Referral{}
	err := r.db.Where("circle_id = ?", circleID).First(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

// FindOneReferralByCode implements ReferralRepo.
func (r *ReferralRepo) FindOneReferralByCode(referralCode string) (*entity.Referral, *domain.Error) {
	referral := &entity.Referral{}
	err := r.db.Where("referral_code = ?", referralCode).First(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

// CreateReferral implements ReferralRepo.
func (r *ReferralRepo) CreateReferral(referral *entity.Referral) (*entity.Referral, *domain.Error) {
	err := r.db.Create(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

func NewReferralRepo(db *gorm.DB) *ReferralRepo {
	return &ReferralRepo{
		db: db,
	}
}
