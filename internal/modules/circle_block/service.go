package circleblock

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"errors"

	"gorm.io/gorm"
)

type CircleBlockService interface {
	CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	DeleteByID(id int) *domain.Error
	GetOneByBlock(prefix string, block string) (*entity.CircleBlock, *domain.Error)
	UpsertOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
}

type circleBlockService struct {
	repo CircleBlockRepo
}

// UpsertOne implements CircleBlockService.
func (c *circleBlockService) UpsertOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	data, err := c.repo.GetOneByBlock(block.Prefix, block.Postfix)
	if err != nil && !errors.Is(err.Err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if data == nil {
		return c.CreateOne(block)
	}
	return c.UpdateOne(data.ID, block)
}

// CreateOne implements CircleBlockService.
func (c *circleBlockService) CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	data, err := c.repo.CreateOne(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteByID implements CircleBlockService.
func (c *circleBlockService) DeleteByID(id int) *domain.Error {
	err := c.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

// GetOneByBlock implements CircleBlockService.
func (c *circleBlockService) GetOneByBlock(prefix string, block string) (*entity.CircleBlock, *domain.Error) {
	data, err := c.repo.GetOneByBlock(prefix, block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateOne implements CircleBlockService.
func (c *circleBlockService) UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	data, err := c.repo.UpdateOne(id, block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewCircleBlockService(repo CircleBlockRepo) CircleBlockService {
	return &circleBlockService{repo: repo}
}
