package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	referral_dto "catalog-be/internal/modules/circle/referral/dto"
)

type ReferralService interface {
	CreateNewReferral(dto referral_dto.CreateReferralBody) (*entity.Referral, *domain.Error)
}

type referralService struct {
	repo ReferralRepo
}

// CreateNewReferral implements ReferralService.
func (r *referralService) CreateNewReferral(dto referral_dto.CreateReferralBody) (*entity.Referral, *domain.Error) {
	created, err := r.repo.CreateReferral(
		&entity.Referral{CircleID: dto.CircleID, ReferralCode: dto.ReferralCode},
	)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func NewReferralService(repo ReferralRepo) ReferralService {
	return &referralService{repo}
}
