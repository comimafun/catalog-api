package circle_work_type

import (
	"gorm.io/gorm"
)

type CircleWorkTypeRepo struct {
	db *gorm.DB
}

func NewCircleWorkTypeRepo(db *gorm.DB) *CircleWorkTypeRepo {
	return &CircleWorkTypeRepo{db: db}
}
