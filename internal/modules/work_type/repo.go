package work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type WorkTypeRepo interface {
	CreateOne(workType entity.WorkType) (*entity.WorkType, *domain.Error)
	UpdateOne(id int, workType entity.WorkType) (*entity.WorkType, *domain.Error)
	DeleteByID(id int) *domain.Error
	FindAll() (*[]entity.WorkType, *domain.Error)
}
type workTypeRepo struct {
	db *gorm.DB
}

// FindAll implements WorkTypeRepo.
func (w *workTypeRepo) FindAll() (*[]entity.WorkType, *domain.Error) {
	var workTypes []entity.WorkType
	if err := w.db.Find(&workTypes).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workTypes, nil
}

// CreateOne implements WorkTypeRepo.
func (w *workTypeRepo) CreateOne(workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	if err := w.db.Create(&workType).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workType, nil
}

// DeleteByID implements WorkTypeRepo.
func (w *workTypeRepo) DeleteByID(id int) *domain.Error {
	if err := w.db.Delete(&entity.WorkType{}, id).Error; err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// UpdateOne implements WorkTypeRepo.
func (w *workTypeRepo) UpdateOne(id int, workType entity.WorkType) (*entity.WorkType, *domain.Error) {
	err := w.db.Where("id = ?", id).Updates(workType).Scan(&workType).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &workType, nil
}

func NewWorkTypeRepo(
	db *gorm.DB,
) WorkTypeRepo {
	return &workTypeRepo{
		db: db,
	}
}
