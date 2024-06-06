package circleblock

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type CircleBlockRepo interface {
	CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error)
	DeleteByID(id int) *domain.Error
	GetOneByBlock(prefix string, block string) (*entity.CircleBlock, *domain.Error)
	GetOneByID(id int) (*entity.CircleBlock, *domain.Error)
}
type circleBlockRepo struct {
	db *gorm.DB
}

// GetOneByID implements CircleBlock.
func (c *circleBlockRepo) GetOneByID(id int) (*entity.CircleBlock, *domain.Error) {
	var circleBlock entity.CircleBlock
	err := c.db.First(&circleBlock, id).Error
	if err != nil {
		return nil, &domain.Error{Err: err, Code: 500}
	}
	return &circleBlock, nil
}

// GetOneByBlock implements CircleBlock.
func (c *circleBlockRepo) GetOneByBlock(prefix string, block string) (*entity.CircleBlock, *domain.Error) {
	var circleBlock entity.CircleBlock
	err := c.db.Where("prefix = ? AND postfix = ?", prefix, block).First(&circleBlock).Error
	if err != nil {
		return nil, &domain.Error{Err: err, Code: 500}
	}
	return &circleBlock, nil
}

// CreateOne implements CircleBlock.
func (c *circleBlockRepo) CreateOne(block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	err := c.db.Create(&block).Error
	if err != nil {
		return nil, &domain.Error{Err: err, Code: 500}
	}
	return &block, nil
}

// DeleteByID implements CircleBlock.
func (c *circleBlockRepo) DeleteByID(id int) *domain.Error {
	err := c.db.Delete(&entity.CircleBlock{}, id).Error
	if err != nil {
		return &domain.Error{Err: err, Code: 500}
	}
	return nil
}

// UpdateOne implements CircleBlock.
func (c *circleBlockRepo) UpdateOne(id int, block entity.CircleBlock) (*entity.CircleBlock, *domain.Error) {
	err := c.db.Model(&entity.CircleBlock{}).Where("id = ?", id).Updates(&block).Scan(&block).Error
	if err != nil {
		return nil, &domain.Error{Err: err, Code: 500}
	}
	return &block, nil
}

func NewCircleBlockRepo(db *gorm.DB) CircleBlockRepo {
	return &circleBlockRepo{db}
}
