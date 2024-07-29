package circle_upvote

import (
	"catalog-be/internal/domain"

	"gorm.io/gorm"
)

type CircleUpvoteRepo struct {
	db *gorm.DB
}

// FindUpvoteByCircleIDAndUserID implements CircleUpvoteRepo.
func (c *CircleUpvoteRepo) FindUpvoteByCircleIDAndUserID(circleID int, userID int) (bool, *domain.Error) {
	var isUpvoted bool
	err := c.db.Raw("SELECT EXISTS(SELECT 1 FROM user_upvote WHERE user_id = ? AND circle_id = ?) as is_upvoted", userID, circleID).Scan(&isUpvoted).Error
	if err != nil {
		return false, domain.NewError(500, err, nil)
	}

	return isUpvoted, nil
}

// CreateUpvote implements CircleUpvoteRepo.
func (c *CircleUpvoteRepo) CreateUpvote(circleID int, userID int) *domain.Error {
	err := c.db.Exec("INSERT INTO user_upvote (user_id, circle_id) VALUES (?, ?)", userID, circleID).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

// DeleteUpvote implements CircleUpvoteRepo.
func (c *CircleUpvoteRepo) DeleteUpvote(circleID int, userID int) *domain.Error {
	err := c.db.Exec("DELETE FROM user_upvote WHERE user_id = ? AND circle_id = ?", userID, circleID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

func NewCircleUpvoteRepo(db *gorm.DB) *CircleUpvoteRepo {
	return &CircleUpvoteRepo{db}
}
