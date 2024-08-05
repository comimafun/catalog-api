package work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"time"
)

type WorkTypeService struct {
	workTypeRepo *WorkTypeRepo
}

// CreateOne implements WorkTypeService.
func (w *WorkTypeService) CreateOne(workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	data, err := w.workTypeRepo.CreateOne(workType)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteByID implements WorkTypeService.
func (w *WorkTypeService) DeleteByID(id int) *domain.Error {
	return w.workTypeRepo.DeleteByID(id)
}

// FindAll implements WorkTypeService.
func (w *WorkTypeService) FindAll() (*[]entity.WorkType, *domain.Error) {
	data, err := w.workTypeRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateOne implements WorkTypeService.
func (w *WorkTypeService) UpdateOne(id int, workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	now := time.Now()
	return w.workTypeRepo.UpdateOne(id, entity.WorkType{
		Name:      workType.Name,
		UpdatedAt: &now,
	})
}

func NewWorkTypeService(
	workTypeRepo *WorkTypeRepo,
) *WorkTypeService {
	return &WorkTypeService{
		workTypeRepo: workTypeRepo,
	}
}
