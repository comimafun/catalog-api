package referral

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	referral_dto "catalog-be/internal/modules/circle/referral/dto"
	"errors"

	"gorm.io/gorm"
)

type ReferralService interface {
	CreateNewReferral(dto referral_dto.CreateReferralBody) (*entity.Referral, *domain.Error)
	FindReferralByCode(referralCode string) (*entity.Referral, *domain.Error)
}

type referralService struct {
	repo ReferralRepo
}

// FindReferralByCode implements ReferralService.
func (r *referralService) FindReferralByCode(referralCode string) (*entity.Referral, *domain.Error) {
	code, err := r.repo.FindOneReferralByCode(referralCode)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, errors.New("NOT_FOUND"), nil)
		}

		return nil, err
	}

	return code, nil
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
