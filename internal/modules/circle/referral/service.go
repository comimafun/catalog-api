package referral

type ReferralService interface{}

type referralService struct {
	repo ReferralRepo
}

func NewReferralService(repo ReferralRepo) ReferralService {
	return &referralService{repo}
}
