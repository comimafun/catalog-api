package circle_fandom

import (
	"gorm.io/gorm"
)

type CircleFandomRepo struct {
	db *gorm.DB
}

func NewCircleFandomRepo(db *gorm.DB) *CircleFandomRepo {
	return &CircleFandomRepo{
		db: db,
	}
}
