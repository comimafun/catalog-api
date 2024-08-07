package entity

import "time"

type Report struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	CircleID  int        `json:"circle_id"`
	Reason    string     `json:"reason"`
	CreatedAt *time.Time `json:"created_at"`
}

func (Report) TableName() string {
	return "report"
}
