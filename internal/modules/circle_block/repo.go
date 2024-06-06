package circleblock

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type CircleBlock interface {
	CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	DeleteByID(id int) *domain.Error
}
type circleBlock struct {
	db *gorm.DB
}

// CreateOne implements CircleBlock.
func (c *circleBlock) CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	err := c.db.Create(&block).Error
	if err != nil {
		return nil, &domain.Error{Err: err}
	}
	return &block, nil
}

// DeleteByID implements CircleBlock.
func (c *circleBlock) DeleteByID(id int) *domain.Error {
	err := c.db.Delete(&entity.CircleBlock{}, id).Error
	if err != nil {
		return &domain.Error{Err: err}
	}
	return nil
}

// UpdateOne implements CircleBlock.
func (c *circleBlock) UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	err := c.db.Model(&entity.CircleBlock{}).Where("id = ?", id).Updates(&block).Scan(&block).Error
	if err != nil {
		return nil, &domain.Error{Err: err}
	}
	return &block, nil
}

func NewCircleBlockRepo(db *gorm.DB) CircleBlock {
	return &circleBlock{db}
}
