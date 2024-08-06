package bookmark

import (
	"catalog-be/internal/domain"
)

type CircleBookmarkService struct {
	circleRepo *CircleBookmarkRepo
}

// DeleteBookmarkByUserCircleID implements CircleBookmarkService.
func (c *CircleBookmarkService) DeleteBookmarkByUserCircleID(circleID int, userID int) *domain.Error {
	return c.circleRepo.DeleteBookmarkByUserCircleID(circleID, userID)
}

// CreateOneBookmark implements CircleBookmarkService.
func (c *CircleBookmarkService) CreateOneBookmark(circleID int, userID int) *domain.Error {
	return c.circleRepo.CreateOneBookmark(circleID, userID)
}

func NewCircleBookmarkService(repo *CircleBookmarkRepo) *CircleBookmarkService {
	return &CircleBookmarkService{circleRepo: repo}
}
