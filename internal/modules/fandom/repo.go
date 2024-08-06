package fandom

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	fandom_dto "catalog-be/internal/modules/fandom/dto"

	"gorm.io/gorm"
)

type FandomRepo struct {
	db *gorm.DB
}

// GetFandomCount implements FandomRepo.
func (f *FandomRepo) GetFandomCount(filter *fandom_dto.GetPaginatedFandomFilter) (int, *domain.Error) {
	var count int64
	err := f.db.Model(&entity.Fandom{}).
		Where("name ilike ? and visible = ?", "%"+filter.Search+"%", true).
		Count(&count).Error
	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}
	return int(count), nil
}

// CreateOneFandom implements FandomRepo.
func (f *FandomRepo) CreateOneFandom(fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	if err := f.db.Create(&fandom).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &fandom, nil
}

// DeleteOneFandomByFandomID implements FandomRepo.
func (f *FandomRepo) DeleteOneFandomByFandomID(id int) *domain.Error {
	if err := f.db.Delete(&entity.Fandom{}, id).Error; err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// GetPaginatedFandoms implements FandomRepo.
func (f *FandomRepo) GetPaginatedFandoms(filter *fandom_dto.GetPaginatedFandomFilter) ([]entity.Fandom, *domain.Error) {
	var fandoms []entity.Fandom
	err := f.db.Where("name ilike ? and visible = ?", "%"+filter.Search+"%", true).
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&fandoms).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return fandoms, nil
}

// UpdateOneFandom implements FandomRepo.
func (f *FandomRepo) UpdateOneFandom(id int, fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	err := f.db.Where("id = ?", id).Updates(fandom).Scan(&fandom).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &fandom, nil
}

func NewFandomRepo(db *gorm.DB) *FandomRepo {
	return &FandomRepo{
		db,
	}
}
