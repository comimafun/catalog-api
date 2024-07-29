package work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type WorkTypeRepo struct {
	db *gorm.DB
}

// FindAll implements WorkTypeRepo.
func (w *WorkTypeRepo) FindAll() (*[]entity.WorkType, *domain.Error) {
	var workTypes []entity.WorkType
	if err := w.db.Find(&workTypes).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workTypes, nil
}

// CreateOne implements WorkTypeRepo.
func (w *WorkTypeRepo) CreateOne(workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	if err := w.db.Create(&workType).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workType, nil
}

// DeleteByID implements WorkTypeRepo.
func (w *WorkTypeRepo) DeleteByID(id int) *domain.Error {
	if err := w.db.Delete(&entity.WorkType{}, id).Error; err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// UpdateOne implements WorkTypeRepo.
func (w *WorkTypeRepo) UpdateOne(id int, workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	err := w.db.Where("id = ?", id).Updates(workType).Scan(&workType).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workType, nil
}

func NewWorkTypeRepo(
	db *gorm.DB,
) *WorkTypeRepo {
	return &WorkTypeRepo{
		db: db,
	}
}
