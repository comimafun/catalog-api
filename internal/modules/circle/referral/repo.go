package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ReferralRepo struct {
	db *gorm.DB
}

// GetOneReferralByCircleID implements ReferralRepo.
func (r *ReferralRepo) GetOneReferralByCircleID(circleID int) (*entity.Referral, *domain.Error) {
	referral := &entity.Referral{}
	err := r.db.Where("circle_id = ?", circleID).First(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

// GetOneReferralByCode implements ReferralRepo.
func (r *ReferralRepo) GetOneReferralByCode(referralCode string) (*entity.Referral, *domain.Error) {
	referral := &entity.Referral{}
	err := r.db.Where("referral_code = ?", referralCode).First(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

// CreateOneReferral implements ReferralRepo.
func (r *ReferralRepo) CreateOneReferral(referral *entity.Referral) (*entity.Referral, *domain.Error) {
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
