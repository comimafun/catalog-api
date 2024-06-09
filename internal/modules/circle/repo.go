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
	FindOneBySlugAndRelatedTables(slug string, userID int) ([]entity.CircleRaw, *domain.Error)

	FindOneByUserID(userID int) (*entity.Circle, *domain.Error)
	UpdateOneByID(circleID int, circle entity.Circle) (*entity.Circle, *domain.Error)
	FindAllCircles(filter *circle_dto.FindAllCircleFilter, userID int) ([]entity.CircleRaw, *domain.Error)
	FindAllCount(filter *circle_dto.FindAllCircleFilter) (int, *domain.Error)
	findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{})
	DeleteOneByID(id int) *domain.Error

	FindAllBookmarkedCount(userID int, filter *circle_dto.FindAllCircleFilter) (int, *domain.Error)
	FindBookmarkedCircleByUserID(userID int, filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error)
}
type circleRepo struct {
	db *gorm.DB
}

func NewCircleRepo(
	db *gorm.DB,
) CircleRepo {
	return &circleRepo{
		db: db,
	}
}

// FindOneBySlugAndRelatedTables implements CircleRepo.
func (c *circleRepo) FindOneBySlugAndRelatedTables(slug string, userID int) ([]entity.CircleRaw, *domain.Error) {
	query := `
		SELECT
			c.*,
			f. "name" AS fandom_name,
			f.id AS fandom_id,
			f.visible AS fandom_visible,
			f.created_at AS fandom_created_at,
			f.updated_at AS fandom_updated_at,
			f.deleted_at AS fandom_updated_at,
			wt. "name" AS work_type_name,
			wt.id AS work_type_id,
			wt.created_at AS work_type_created_at,
			wt.updated_at AS work_type_updated_at,
			wt.deleted_at AS work_type_updated_at,
			p.id as product_id,
			p."name" as product_name,
			p.image_url as product_image_url,
			p.created_at as product_created_at,
			p.updated_at as product_updated_at,
			cb.id AS block_id,
			cb.prefix AS block_prefix,
			cb.postfix AS block_postfix,
			cb.created_at AS block_created_at,
			cb.updated_at AS block_updated_at,
			ub.created_at AS bookmarked_at,
			CASE WHEN ub.user_id IS NOT NULL THEN
				TRUE
			ELSE
				FALSE
			END AS bookmarked
		FROM
			circle c
		LEFT JOIN circle_fandom cf ON c.id = cf.circle_id
		LEFT JOIN fandom f ON f.id = cf.fandom_id
		LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id
		LEFT JOIN work_type wt ON wt.id = cwt.work_type_id
		LEFT JOIN circle_block cb ON c.id = cb.id
		LEFT JOIN product p on c.id = p.circle_id
		LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)
		WHERE
			c.deleted_at IS NULL
			AND c.slug = ?
		`

	var row []entity.CircleRaw

	err := c.db.Raw(query, userID, slug).Scan(&row).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return row, nil
}

// findAllBookmarkedCount implements CircleRepo.
func (c *circleRepo) FindAllBookmarkedCount(userID int, filter *circle_dto.FindAllCircleFilter) (int, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	query := fmt.Sprintf(`
		select
			count(DISTINCT c.id) as count
		from 
			circle c
		JOIN
			user_bookmark ub on c.id = ub.circle_id and ub.user_id = ?
		LEFT JOIN
			circle_fandom cf on c.id = cf.circle_id
		LEFT JOIN
			fandom f on f.id = cf.fandom_id
		lEFT JOIN
			circle_work_type cwt on c.id = cwt.circle_id
		LEFT JOIN
			work_type wt on wt.id = cwt.work_type_id
		%s
	`, whereClause)

	args = append(args, userID)

	var count int
	err := c.db.Raw(query, args...).Scan(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return count, nil
}

// FindBookmarkedCircleByUserID implements CircleRepo.
func (c *circleRepo) FindBookmarkedCircleByUserID(userID int, filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	query := fmt.Sprintf(`
		select
			c.*,
			
			f."name" as fandom_name,
			f.id as fandom_id,
			f.visible as fandom_visible,
			f.created_at as fandom_created_at,
			f.updated_at as fandom_updated_at,
			f.deleted_at as fandom_updated_at,

			wt."name" as work_type_name,
			wt.id as work_type_id,
			wt.created_at as work_type_created_at,
			wt.updated_at as work_type_updated_at,
			wt.deleted_at as work_type_updated_at,

			cb.id as block_id,
			cb.prefix as block_prefix,
			cb.postfix as block_postfix,
			cb.created_at as block_created_at,
			cb.updated_at as block_updated_at,

			ub.created_at as bookmarked_at,
			CASE WHEN ub.user_id is not null THEN true ELSE false END as bookmarked
		from 
			circle c
		JOIN
			user_bookmark ub on c.id = ub.circle_id and ub.user_id = ?
		LEFT JOIN
			circle_fandom cf on c.id = cf.circle_id
		LEFT JOIN
			fandom f on f.id = cf.fandom_id
		lEFT JOIN
			circle_work_type cwt on c.id = cwt.circle_id
		LEFT JOIN
			work_type wt on wt.id = cwt.work_type_id
		LEFT JOIN
			circle_block cb on c.id = cb.circle_id
		%s
		order by ub.created_at desc
		offset ?
		limit ?
	`, whereClause)

	args = append(args, userID)

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, offset, filter.Limit)

	var circles []entity.CircleRaw
	err := c.db.Raw(query, args...).Scan(&circles).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circles, nil
}

// findAllWhereSQL implements CircleRepo.
func (c *circleRepo) findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{}) {
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

	if len(filter.WorkTypeID) > 0 {
		whereClause += " and wt.id in (?)"
		args = append(args, filter.WorkTypeID)
	}

	return whereClause, args
}

// FindAll implements CircleRepo.
func (c *circleRepo) FindAllCircles(filter *circle_dto.FindAllCircleFilter, userID int) ([]entity.CircleRaw, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	query := fmt.Sprintf(`
		SELECT
			c.*,
			
			f."name" as fandom_name,
			f.id as fandom_id,
			f.visible as fandom_visible,
			f.created_at as fandom_created_at,
			f.updated_at as fandom_updated_at,
			f.deleted_at as fandom_updated_at,

			wt."name" as work_type_name,
			wt.id as work_type_id,
			wt.created_at as work_type_created_at,
			wt.updated_at as work_type_updated_at,
			wt.deleted_at as work_type_updated_at,

			cb.id as block_id,
			cb.prefix as block_prefix,
			cb.postfix as block_postfix,
			cb.created_at as block_created_at,
			cb.updated_at as block_updated_at,

			ub.created_at as bookmarked_at,
			CASE WHEN ub.user_id is not null THEN true ELSE false END as bookmarked
		FROM 
			circle c
		LEFT JOIN
			circle_fandom cf on c.id = cf.circle_id
		LEFT JOIN
			fandom f on f.id = cf.fandom_id
		lEFT JOIN
			circle_work_type cwt on c.id = cwt.circle_id
		LEFT JOIN
			work_type wt on wt.id = cwt.work_type_id
		LEFT JOIN
			user_bookmark ub on c.id = ub.circle_id and ub.user_id = COALESCE(?, ub.user_id)
		LEFT JOIN
			circle_block cb on c.id = cb.circle_id
		%s
		ORDER BY
			c.created_at desc
		OFFSET ?
		LIMIT ?
	`, whereClause)

	offset := (filter.Page - 1) * filter.Limit
	args = append(args, userID)

	args = append(args, offset, filter.Limit)

	var circles []entity.CircleRaw
	err := c.db.Raw(query, args...).Scan(&circles).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circles, nil
}

// FindAllCount implements CircleRepo.
func (c *circleRepo) FindAllCount(filter *circle_dto.FindAllCircleFilter) (int, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	query := fmt.Sprintf(`
		select
			count(DISTINCT c.id) as count
		from 
			circle c
		LEFT JOIN
			circle_fandom cf on c.id = cf.circle_id
		LEFT JOIN
			fandom f on f.id = cf.fandom_id
		lEFT JOIN
			circle_work_type cwt on c.id = cwt.circle_id
		LEFT JOIN
			work_type wt on wt.id = cwt.work_type_id
		%s
	`, whereClause)

	var count int
	err := c.db.Raw(query, args...).Scan(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return count, nil
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
