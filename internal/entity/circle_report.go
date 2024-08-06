package entity

import "time"

type CircleReport struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	CircleID  int        `json:"circle_id"`
	Reason    string     `json:"reason"`
	CreatedAt *time.Time `json:"created_at"`
}

func (CircleReport) TableName() string {
	return "circle_report"
}
