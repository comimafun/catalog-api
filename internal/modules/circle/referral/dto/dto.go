package referral_dto

type CreateReferralPayload struct {
	ReferralCode string `json:"referral_code" validate:"required"`
	CircleID     int    `json:"circle_id" validate:"required"`
}
