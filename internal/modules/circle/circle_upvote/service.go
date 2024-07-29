package circle_upvote

import "catalog-be/internal/domain"

type CircleUpvoteService struct {
	repo *CircleUpvoteRepo
}

// CancelUpvoteCircle implements CircleUpvoteService.
func (c *CircleUpvoteService) CancelUpvoteCircle(circleID int, userID int) *domain.Error {
	return c.repo.DeleteUpvote(circleID, userID)
}

// IsUpvoted implements CircleUpvoteService.
func (c *CircleUpvoteService) IsUpvoted(circleID int, userID int) (bool, *domain.Error) {
	return c.repo.FindUpvoteByCircleIDAndUserID(circleID, userID)
}

// UpvoteCircle implements CircleUpvoteService.
func (c *CircleUpvoteService) UpvoteCircle(circleID int, userID int) *domain.Error {
	return c.repo.CreateUpvote(circleID, userID)
}

func NewCircleUpvoteService(repo *CircleUpvoteRepo) *CircleUpvoteService {
	return &CircleUpvoteService{repo}
}
