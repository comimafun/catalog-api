package circle_work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleWorkTypeService struct {
	repo *CircleWorkTypeRepo
}

// BulkDeleteAndInsertCircleWorkType implements CircleWorkTypeService.
func (c *CircleWorkTypeService) BulkDeleteAndInsertCircleWorkType(circleID int, workTypeIDs []int) *domain.Error {
	return c.repo.BatchInsertCircleWorkTypeRelation(circleID, workTypeIDs)
}

// FindAllCircleWorkTypeRelated implements CircleWorkTypeService.
func (c *CircleWorkTypeService) FindAllCircleWorkTypeRelated(circleID int) ([]entity.WorkType, *domain.Error) {
	return c.repo.FindAllCircleRelationWorkType(circleID)
}

func NewCircleWorkTypeService(repo *CircleWorkTypeRepo) *CircleWorkTypeService {
	return &CircleWorkTypeService{repo: repo}
}
