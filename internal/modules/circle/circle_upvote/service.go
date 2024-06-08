package circle_upvote

import "catalog-be/internal/domain"

type CircleUpvoteService interface {
	UpvoteCircle(circleID int, userID int) *domain.Error
	CancelUpvoteCircle(circleID int, userID int) *domain.Error
	IsUpvoted(circleID int, userID int) (bool, *domain.Error)
}

type circleUpvoteService struct {
	repo CircleUpvoteRepo
}

// CancelUpvoteCircle implements CircleUpvoteService.
func (c *circleUpvoteService) CancelUpvoteCircle(circleID int, userID int) *domain.Error {
	return c.repo.DeleteUpvote(circleID, userID)
}

// IsUpvoted implements CircleUpvoteService.
func (c *circleUpvoteService) IsUpvoted(circleID int, userID int) (bool, *domain.Error) {
	return c.repo.FindUpvoteByCircleIDAndUserID(circleID, userID)
}

// UpvoteCircle implements CircleUpvoteService.
func (c *circleUpvoteService) UpvoteCircle(circleID int, userID int) *domain.Error {
	return c.repo.CreateUpvote(circleID, userID)
}

func NewCircleUpvoteService(repo CircleUpvoteRepo) CircleUpvoteService {
	return &circleUpvoteService{repo}
}
