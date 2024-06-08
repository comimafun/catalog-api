package circle

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CircleRepo interface {
	CreateOne(circle entity.Circle) (*entity.Circle, *domain.Error)
	FindOneByID(id int) (*entity.Circle, *domain.Error)
	FindOneBySlug(slug string) (*entity.Circle, *domain.Error)
	FindOneByUserID(userID int) (*entity.Circle, *domain.Error)
	UpdateOneByID(circleID int, circle entity.Circle) (*entity.Circle, *domain.Error)
	FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error)
	FindAllCount(filter *circle_dto.FindAllCircleFilter) (int, *domain.Error)
	findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{})
	DeleteOneByID(id int) *domain.Error
	DeleteWorkTypeRelationByCircleID(circleID int) *domain.Error
	UpsertCircleFandomRelation(circleID int, fandomID int) *domain.Error
	UpsertCircleWorkTypeRelation(circleID int, workTypeID int) *domain.Error
	BatchInsertCircleWorkTypeRelation(circleID int, workTypeIDs []int) *domain.Error
	FindAllCircleRelationFandom(circleID int) ([]entity.Fandom, *domain.Error)
	FindAllCircleRelationWorkType(circleID int) ([]entity.WorkType, *domain.Error)
}
type circleRepo struct {
	db *gorm.DB
}

// BatchInsertCircleWorkTypeRelation implements CircleRepo.
func (c *circleRepo) BatchInsertCircleWorkTypeRelation(circleID int, workTypeIDs []int) *domain.Error {
	var valueStrings []string
	valueArgs := make([]interface{}, 0)
	for _, workTypeID := range workTypeIDs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, circleID, workTypeID)
	}

	query := fmt.Sprintf(`
			INSERT INTO circle_work_type (circle_id, work_type_id)
			VALUES %s
		`, strings.Join(valueStrings, ", "))

	err := c.db.Exec(query, valueArgs...).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// DeleteWorkTypeRelationByCircleID implements CircleRepo.
func (c *circleRepo) DeleteWorkTypeRelationByCircleID(circleID int) *domain.Error {
	err := c.db.Exec(`
		delete from circle_work_type where circle_id = ?
	`, circleID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// UpsertCircleWorkTypeRelation implements CircleRepo.
func (c *circleRepo) UpsertCircleWorkTypeRelation(circleID int, workTypeID int) *domain.Error {
	err := c.db.Exec(`
		insert into circle_work_type (circle_id, work_type_id)
		values (?, ?)
		on conflict (circle_id, work_type_id) do nothing
	`, circleID, workTypeID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindAllCircleRelationWorkType implements CircleRepo.
func (c *circleRepo) FindAllCircleRelationWorkType(circleID int) ([]entity.WorkType, *domain.Error) {
	var workTypes []entity.WorkType
	err := c.db.Raw(`
		select
			w.*
		from
			work_type w
		inner join circle_work_type cw on w.id = cw.work_type_id
		where
			cw.circle_id = ? and w.deleted_at is null
	`, circleID).Scan(&workTypes).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return workTypes, nil
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
func (c *circleRepo) FindAll(filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error) {
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
			wt.deleted_at as work_type_updated_at
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
