package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	referral_dto "catalog-be/internal/modules/circle/referral/dto"
	"errors"

	"gorm.io/gorm"
)

type ReferralService struct {
	repo *ReferralRepo
}

func NewReferralService(repo *ReferralRepo) *ReferralService {
	return &ReferralService{repo}
}

// GetOneReferralCodeByCircleID implements ReferralService.
func (r *ReferralService) GetOneReferralCodeByCircleID(circleID int) (*entity.Referral, *domain.Error) {
	ref, err := r.repo.GetOneReferralByCircleID(circleID)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, errors.New("NOT_FOUND"), nil)
		}

		return nil, err
	}

	return ref, nil
}

// GetOneReferralByCode implements ReferralService.
func (r *ReferralService) GetOneReferralByCode(referralCode string) (*entity.Referral, *domain.Error) {
	code, err := r.repo.GetOneReferralByCode(referralCode)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, errors.New("NOT_FOUND"), nil)
		}

		return nil, err
	}

	return code, nil
}

// CreateOneReferral implements ReferralService.
func (r *ReferralService) CreateOneReferral(dto referral_dto.CreateReferralPayload) (*entity.Referral, *domain.Error) {
	created, err := r.repo.CreateOneReferral(
		&entity.Referral{CircleID: dto.CircleID, ReferralCode: dto.ReferralCode},
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}
