package fandom

import (
	"catalog-be/internal/database/factory"
	"catalog-be/internal/domain"
	"catalog-be/internal/dto"
	"catalog-be/internal/entity"
	fandom_dto "catalog-be/internal/modules/fandom/dto"
	"time"
)

type FandomService struct {
	fandomRepo *FandomRepo
}

// DeleteOneFandomByFandomID implements FandomService.
func (f *FandomService) DeleteOneFandomByFandomID(id int) *domain.Error {
	return f.fandomRepo.DeleteOneFandomByFandomID(id)
}

// GetPaginatedFandoms implements FandomService.
func (f *FandomService) GetPaginatedFandoms(filter *fandom_dto.GetPaginatedFandomFilter) (*dto.Pagination[[]entity.Fandom], *domain.Error) {
	count, countErr := f.fandomRepo.GetFandomCount(filter)
	if countErr != nil {
		return nil, countErr
	}
	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	fandoms, findErr := f.fandomRepo.GetPaginatedFandoms(filter)
	if findErr != nil {
		return nil, findErr
	}
	return &dto.Pagination[[]entity.Fandom]{
		Data:     fandoms,
		Metadata: *metadata,
	}, nil
}

// UpdateOneFandomByFandomID implements FandomService.
func (f *FandomService) UpdateOneFandomByFandomID(id int, fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	now := time.Now()
	return f.fandomRepo.UpdateOneFandom(id, entity.Fandom{
		Name:      fandom.Name,
		Visible:   fandom.Visible,
		UpdatedAt: &now,
	})
}

// CreateOneFandom implements FandomService.
func (f *FandomService) CreateOneFandom(fandom entity.Fandom) (*entity.Fandom, *domain.Error) {
	return f.fandomRepo.CreateOneFandom(fandom)
}

func NewFandomService(fandomRepo *FandomRepo) *FandomService {
	return &FandomService{
		fandomRepo,
	}
}
