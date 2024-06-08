package circle

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"fmt"

	"gorm.io/gorm"
)

type CircleRepo interface {
	CreateOne(circle entity.Circle) (*entity.Circle, *domain.Error)
	FindOneByID(id int) (*entity.Circle, *domain.Error)
	FindOneBySlug(slug string) (*entity.Circle, *domain.Error)
	FindOneByUserID(userID int) (*entity.Circle, *domain.Error)
	UpdateOneByID(circleID int, circle entity.Circle) (*entity.Circle, *domain.Error)
	FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.Circle, *domain.Error)
	DeleteOneByID(id int) *domain.Error
}
type circleRepo struct {
	db *gorm.DB
}

// FindAll implements CircleRepo.
func (c *circleRepo) FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.Circle, *domain.Error) {
	whereClause := "where c.deleted_at is null"
	args := make([]interface{}, 0)

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" and c.name ilike '%%?%%'")
		args = append(args, filter.Search)
	}

	query := fmt.Sprintf(`
		select
			c.*
		from
			circle c
		%s
		offset ?
		limit ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, offset, filter.Limit)

	var circles []entity.Circle
	err := c.db.Raw(query, args...).Scan(&circles).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return circles, nil
}

// CreateOne implements CircleRepo.
func (c *circleRepo) CreateOne(circle entity.Circle) (*entity.Circle, *domain.Error) {
	err := c.db.Create(&circle).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &circle, nil
}

// DeleteOneByID implements CircleRepo.
func (c *circleRepo) DeleteOneByID(id int) *domain.Error {
	err := c.db.Delete(&entity.Circle{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindOneByID implements CircleRepo.
func (c *circleRepo) FindOneByID(id int) (*entity.Circle, *domain.Error) {
	var circle entity.Circle
	err := c.db.First(&circle, id).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &circle, nil
}

// FindOneBySlug implements CircleRepo.
func (c *circleRepo) FindOneBySlug(slug string) (*entity.Circle, *domain.Error) {
	var circle entity.Circle
	err := c.db.First(&circle, "slug = ?", slug).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &circle, nil
}

// FindOneByUserID implements CircleRepo.
func (c *circleRepo) FindOneByUserID(userID int) (*entity.Circle, *domain.Error) {
	var circle entity.Circle
	err := c.db.First(&circle, "user_id = ?", userID).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &circle, nil
}

// UpdateOneByID implements CircleRepo.
func (c *circleRepo) UpdateOneByID(circleID int, circle entity.Circle) (*entity.Circle, *domain.Error) {
	var updated entity.Circle
	err := c.db.Table("circle").Where("id = ?", circleID).Updates(&circle).Scan(&updated).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &updated, nil
}

func NewCircleRepo(
	db *gorm.DB,
) CircleRepo {
	return &circleRepo{
		db: db,
	}
}
