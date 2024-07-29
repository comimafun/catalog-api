package bookmark

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleBookmarkService struct {
	circleRepo *CircleBookmarkRepo
}

// FindByCircleIDAndUserID implements CircleBookmarkService.
func (c *CircleBookmarkService) FindByCircleIDAndUserID(circleID int, userID int) (*entity.UserBookmark, *domain.Error) {
	return c.circleRepo.FindByCircleIDAndUserID(circleID, userID)
}

// DeleteBookmark implements CircleBookmarkService.
func (c *CircleBookmarkService) DeleteBookmark(circleID int, userID int) *domain.Error {
	return c.circleRepo.DeleteBookmark(circleID, userID)
}

// CreateBookmark implements CircleBookmarkService.
func (c *CircleBookmarkService) CreateBookmark(circleID int, userID int) *domain.Error {
	return c.circleRepo.CreateBookmark(circleID, userID)
}

func NewCircleBookmarkService(repo *CircleBookmarkRepo) *CircleBookmarkService {
	return &CircleBookmarkService{circleRepo: repo}
}
