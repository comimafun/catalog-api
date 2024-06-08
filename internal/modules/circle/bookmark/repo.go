package bookmark

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type CircleBookmarkRepo interface {
	CreateBookmark(circleID int, userID int) *domain.Error
	DeleteBookmark(circleID int, userID int) *domain.Error
}

type circleBookmarkRepo struct {
	db *gorm.DB
}

// DeleteBookmark implements CircleBookmarkRepo.
func (c *circleBookmarkRepo) DeleteBookmark(circleID int, userID int) *domain.Error {
	err := c.db.Table("user_bookmark").Where("circle_id = ? AND user_id = ?", circleID, userID).Delete(&entity.UserBookmark{}).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

// CreateBookmark implements CircleBookmarkRepo.
func (c *circleBookmarkRepo) CreateBookmark(circleID int, userID int) *domain.Error {
	err := c.db.Table("user_bookmark").Create(&entity.UserBookmark{
		UserID:   userID,
		CircleID: circleID,
	}).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

func NewCircleBookmarkRepo(db *gorm.DB) CircleBookmarkRepo {
	return &circleBookmarkRepo{db: db}
}
