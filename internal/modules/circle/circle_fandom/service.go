package circle_fandom

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleFandomService interface {
	FindAllCircleFandomTypeRelated(circleID int) ([]entity.Fandom, *domain.Error)
	BulkDeleteAndInsertCircleFandomType(circleID int, fandomIDs []int) *domain.Error
}

type circleFandomService struct {
	repo CircleFandomRepo
}

// BulkDeleteAndInsertCircleFandomType implements CircleFandomService.
func (c *circleFandomService) BulkDeleteAndInsertCircleFandomType(circleID int, fandomIDs []int) *domain.Error {
	return c.repo.BatchInsertFandomCircleRelation(circleID, fandomIDs)
}

// FindAllCircleFandomTypeRelated implements CircleFandomService.
func (c *circleFandomService) FindAllCircleFandomTypeRelated(circleID int) ([]entity.Fandom, *domain.Error) {
	return c.repo.FindAllCircleRelationFandom(circleID)
}

func NewCircleFandomService(repo CircleFandomRepo) CircleFandomService {
	return &circleFandomService{repo: repo}
}
