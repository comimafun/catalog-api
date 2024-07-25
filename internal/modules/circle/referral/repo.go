package referral

import "gorm.io/gorm"

type ReferralRepo interface{}

type referralRepo struct {
	db *gorm.DB
}

func NewReferralRepo(db *gorm.DB) ReferralRepo {
	return &referralRepo{db}
}
