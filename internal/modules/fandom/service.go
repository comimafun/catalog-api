package fandom

import (
	"catalog-be/internal/database/factory"
	"catalog-be/internal/domain"
	"catalog-be/internal/dto"
	"catalog-be/internal/entity"
	fandom_dto "catalog-be/internal/modules/fandom/dto"
	"time"
)

type FandomService interface {
	CreateOne(fandom entity.Fandom) (*entity.Fandom, *domain.Error)
	UpdateOne(id int, fandom entity.Fandom) (*entity.Fandom, *domain.Error)
	DeleteByID(id int) *domain.Error
	GetPaginatedFandom(filter *fandom_dto.FindAllFilter) (*dto.Pagination[[]entity.Fandom], *domain.Error)
}

type fandomService struct {
	fandomRepo FandomRepo
}

// DeleteByID implements FandomService.
func (f *fandomService) DeleteByID(id int) *domain.Error {
	return f.fandomRepo.DeleteByID(id)
}

// GetPaginatedFandom implements FandomService.
func (f *fandomService) GetPaginatedFandom(filter *fandom_dto.FindAllFilter) (*dto.Pagination[[]entity.Fandom], *domain.Error) {
	count, countErr := f.fandomRepo.GetFandomCount(filter)
	if countErr != nil {
		return nil, countErr
	}
	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	fandoms, findErr := f.fandomRepo.FindAll(filter)
	if findErr != nil {
		return nil, findErr
	}
	return &dto.Pagination[[]entity.Fandom]{
		Data:     fandoms,
		Metadata: *metadata,
	}, nil
}

// UpdateOne implements FandomService.
func (f *fandomService) UpdateOne(id int, fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	now := time.Now()
	return f.fandomRepo.UpdateOne(id, entity.Fandom{
		Name:      fandom.Name,
		Visible:   fandom.Visible,
		UpdatedAt: &now,
	})
}

// CreateOne implements FandomService.
func (f *fandomService) CreateOne(fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	return f.fandomRepo.CreateOne(fandom)
}

func NewFandomService(fandomRepo FandomRepo) FandomService {
	return &fandomService{
		fandomRepo,
	}
}
