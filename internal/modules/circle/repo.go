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
	FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error)
	DeleteOneByID(id int) *domain.Error
	UpsertCircleFandomRelation(circleID int, fandomID int) *domain.Error
	FindAllCircleRelationFandom(circleID int) ([]entity.Fandom, *domain.Error)
}
type circleRepo struct {
	db *gorm.DB
}

// FindAllCircleRelationFandom implements CircleRepo.
func (c *circleRepo) FindAllCircleRelationFandom(circleID int) ([]entity.Fandom, *domain.Error) {
	var fandoms []entity.Fandom
	err := c.db.Raw(`
		select
			f.*
		from
			fandom f
		inner join circle_fandom cf on f.id = cf.fandom_id
		where
			cf.circle_id = ? and f.deleted_at is null
	`, circleID).Scan(&fandoms).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return fandoms, nil
}

// UpsertCircleFandomRelation implements CircleRepo.
func (c *circleRepo) UpsertCircleFandomRelation(circleID int, fandomID int) *domain.Error {
	err := c.db.Exec(`
		insert into circle_fandom (circle_id, fandom_id)
		values (?, ?)
		on conflict (circle_id, fandom_id) do nothing
	`, circleID, fandomID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindAll implements CircleRepo.
func (c *circleRepo) FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error) {
	whereClause := "where c.deleted_at is null"
	args := make([]interface{}, 0)

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" and c.name ilike '%%?%%'")
		args = append(args, filter.Search)
	}

	if len(filter.FandomID) > 0 {
		whereClause += " and f.id in (?)"
		args = append(args, filter.FandomID)
	}

	query := fmt.Sprintf(`
		select
			c.*,
			
			f."name" as fandom_name,
			f.id as fandom_id,
			f.visible as fandom_visible,
			f.created_at as fandom_created_at,
			f.updated_at as fandom_updated_at,
			f.deleted_at as fandom_updated_at
		from 
			circle c
		LEFT JOIN
			circle_fandom cf on c.id = cf.circle_id
		LEFT JOIN
			fandom f on f.id = cf.fandom_id
		%s
		order by c.created_at desc
		offset ?
		limit ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, offset, filter.Limit)

	var circles []entity.CircleRaw
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
