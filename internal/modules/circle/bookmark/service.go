package bookmark

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleBookmarkService interface {
	CreateBookmark(circleID int, userID int) *domain.Error
	DeleteBookmark(circleID int, userID int) *domain.Error
	FindByCircleIDAndUserID(circleID int, userID int) (*entity.UserBookmark, *domain.Error)
}

type circleBookmarkService struct {
	circleRepo CircleBookmarkRepo
}

// FindByCircleIDAndUserID implements CircleBookmarkService.
func (c *circleBookmarkService) FindByCircleIDAndUserID(circleID int, userID int) (*entity.UserBookmark, *domain.Error) {
	return c.circleRepo.FindByCircleIDAndUserID(circleID, userID)
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
