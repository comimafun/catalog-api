package entity

import "time"

type Referral struct {
	ID           int        `json:"id"`
	ReferralCode string     `json:"referral_code"`
	CircleID     int        `json:"circle_id"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func (Referral) TableName() string {
	return "referral"
}
