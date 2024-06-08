package bookmark

import "catalog-be/internal/domain"

type CircleBookmarkService interface {
	CreateBookmark(circleID int, userID int) *domain.Error
	DeleteBookmark(circleID int, userID int) *domain.Error
}

type circleBookmarkService struct {
	circleRepo CircleBookmarkRepo
}

// DeleteBookmark implements CircleBookmarkService.
func (c *circleBookmarkService) DeleteBookmark(circleID int, userID int) *domain.Error {
	return c.circleRepo.DeleteBookmark(circleID, userID)
}

// CreateBookmark implements CircleBookmarkService.
func (c *circleBookmarkService) CreateBookmark(circleID int, userID int) *domain.Error {
	return c.circleRepo.CreateBookmark(circleID, userID)
}

func NewCircleBookmarkService(repo CircleBookmarkRepo) CircleBookmarkService {
	return &circleBookmarkService{circleRepo: repo}
}
