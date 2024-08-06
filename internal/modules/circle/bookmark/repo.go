package bookmark

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type CircleBookmarkRepo struct {
	db *gorm.DB
}

// DeleteBookmarkByUserCircleID implements CircleBookmarkRepo.
func (c *CircleBookmarkRepo) DeleteBookmarkByUserCircleID(circleID int, userID int) *domain.Error {
	err := c.db.Table("user_bookmark").Where("circle_id = ? AND user_id = ?", circleID, userID).Delete(&entity.UserBookmark{}).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

// CreateOneBookmark implements CircleBookmarkRepo.
func (c *CircleBookmarkRepo) CreateOneBookmark(circleID int, userID int) *domain.Error {
	err := c.db.Table("user_bookmark").Create(&entity.UserBookmark{
		UserID:   userID,
		CircleID: circleID,
	}).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

func NewCircleBookmarkRepo(db *gorm.DB) *CircleBookmarkRepo {
	return &CircleBookmarkRepo{db: db}
}
