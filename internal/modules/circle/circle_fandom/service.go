package circle_fandom

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleFandomService struct {
	repo *CircleFandomRepo
}

// BulkDeleteAndInsertCircleFandomType implements CircleFandomService.
func (c *CircleFandomService) BulkDeleteAndInsertCircleFandomType(circleID int, fandomIDs []int) *domain.Error {
	return c.repo.BatchInsertFandomCircleRelation(circleID, fandomIDs)
}

// FindAllCircleFandomTypeRelated implements CircleFandomService.
func (c *CircleFandomService) FindAllCircleFandomTypeRelated(circleID int) ([]entity.Fandom, *domain.Error) {
	return c.repo.FindAllCircleRelationFandom(circleID)
}

func NewCircleFandomService(repo *CircleFandomRepo) *CircleFandomService {
	return &CircleFandomService{repo: repo}
}
