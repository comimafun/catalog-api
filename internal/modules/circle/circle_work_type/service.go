package circle_work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleWorkTypeService interface {
	FindAllCircleWorkTypeRelated(circleID int) ([]entity.WorkType, *domain.Error)
	BulkDeleteAndInsertCircleWorkType(circleID int, workTypeIDs []int) *domain.Error
}

type circleWorkTypeService struct {
	repo CircleWorkTypeRepo
}

// BulkDeleteAndInsertCircleWorkType implements CircleWorkTypeService.
func (c *circleWorkTypeService) BulkDeleteAndInsertCircleWorkType(circleID int, workTypeIDs []int) *domain.Error {
	return c.repo.BatchInsertCircleWorkTypeRelation(circleID, workTypeIDs)
}

// FindAllCircleWorkTypeRelated implements CircleWorkTypeService.
func (c *circleWorkTypeService) FindAllCircleWorkTypeRelated(circleID int) ([]entity.WorkType, *domain.Error) {
	return c.repo.FindAllCircleRelationWorkType(circleID)
}

func NewCircleWorkTypeService(repo CircleWorkTypeRepo) CircleWorkTypeService {
	return &circleWorkTypeService{repo: repo}
}
