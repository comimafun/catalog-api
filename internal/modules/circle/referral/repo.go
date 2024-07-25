package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ReferralRepo interface {
	CreateReferral(referral *entity.Referral) (*entity.Referral, *domain.Error)
	FindOneReferralByCode(referralCode string) (*entity.Referral, *domain.Error)
}

type referralRepo struct {
	db *gorm.DB
}

// FindOneReferralByCode implements ReferralRepo.
func (r *referralRepo) FindOneReferralByCode(referralCode string) (*entity.Referral, *domain.Error) {
	referral := &entity.Referral{}
	err := r.db.Where("referral_code = ?", referralCode).First(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

// CreateReferral implements ReferralRepo.
func (r *referralRepo) CreateReferral(referral *entity.Referral) (*entity.Referral, *domain.Error) {
	err := r.db.Create(referral).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return referral, nil
}

func NewReferralRepo(db *gorm.DB) ReferralRepo {
	return &referralRepo{db}
}
