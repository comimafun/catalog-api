package circle

import (
	"catalog-be/internal/domain"
	"errors"
	"strings"

	"catalog-be/internal/entity"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"fmt"

	"gorm.io/gorm"
)

type CircleRepo interface {
	CreateOne(circle entity.Circle) (*entity.Circle, *domain.Error)
	OnboardNewCircle(circle *entity.Circle, user *entity.User) (*entity.Circle, *domain.Error)
	FindOneByID(id int) (*entity.Circle, *domain.Error)
	FindOneBySlugAndRelatedTables(slug string, userID int) ([]entity.CircleRaw, *domain.Error)

	UpserstOneCircle(circle *entity.Circle) (*entity.Circle, *domain.Error)
	UpdateCircleAndAllRelation(userID int, payload *entity.Circle, body *circle_dto.UpdateCircleRequestBody) ([]entity.CircleRaw, *domain.Error)

	ResetAttendingEvent(circle *entity.Circle) *domain.Error
	UpdateAttendingEvent(circle *entity.Circle, body *circle_dto.UpdateCircleAttendingEvent) *domain.Error

	transformBlockStringIntoBlockEvent(block string) (*entity.BlockEvent, *domain.Error)

	FindAllCircles(filter *circle_dto.FindAllCircleFilter, userID int) ([]entity.CircleRaw, *domain.Error)

	FindAllCount(filter *circle_dto.FindAllCircleFilter) (int, *domain.Error)
	findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{})

	FindAllBookmarkedCount(userID int, filter *circle_dto.FindAllCircleFilter) (int, *domain.Error)
	FindBookmarkedCircleByUserID(userID int, filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error)
}
type circleRepo struct {
	db *gorm.DB
}

// UpdateAttendingEvent implements CircleRepo.
func (c *circleRepo) UpdateAttendingEvent(circle *entity.Circle, body *circle_dto.UpdateCircleAttendingEvent) *domain.Error {
	tx := c.db.Begin()
	if tx.Error != nil {
		return domain.NewError(500, tx.Error, nil)
	}

	err := tx.Save(circle).Error
	if err != nil {
		tx.Rollback()
		return domain.NewError(500, err, nil)
	}

	if body.CircleBlock != "" {
		block, err := c.transformBlockStringIntoBlockEvent(body.CircleBlock)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingBlock := new(entity.BlockEvent)
		existingErr := tx.Where("prefix = ? AND postfix = ?  AND event_id = ?", block.Prefix, block.Postfix, body.EventID).First(&existingBlock).Error

		if existingErr != nil && !errors.Is(existingErr, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return domain.NewError(500, existingErr, nil)
		}

		if existingBlock.ID == 0 {
			// delete previous block
			deleteErr := tx.Table("block_event").Where("circle_id = ? AND event_id = ?", circle.ID, body.EventID).Unscoped().Delete(&entity.BlockEvent{}).Error
			if deleteErr != nil {
				tx.Rollback()
				return domain.NewError(500, deleteErr, nil)
			}

			block.CircleID = circle.ID
			block.EventID = body.EventID

			createErr := tx.Create(block).Error
			if createErr != nil {
				tx.Rollback()
				return domain.NewError(500, createErr, nil)
			}
		} else if existingBlock.CircleID != circle.ID {
			tx.Rollback()
			return domain.NewError(400, errors.New("BLOCK_ALREADY_EXIST"), nil)
		}

	} else {
		// delete block
		err := tx.Table("block_event").Where("circle_id = ? AND event_id = ?", circle.ID, body.EventID).Unscoped().Delete(&entity.BlockEvent{}).Error
		if err != nil {
			tx.Rollback()
			return domain.NewError(500, err, nil)
		}
	}

	tx.Commit()

	return nil
}

// ResetAttendingEvent implements CircleRepo.
func (c *circleRepo) ResetAttendingEvent(circle *entity.Circle) *domain.Error {
	tx := c.db.Begin()
	if tx.Error != nil {
		return domain.NewError(500, tx.Error, nil)
	}

	err := tx.Table("block_event").
		Where("circle_id = ? AND event_id = ?", circle.ID, circle.EventID).
		Unscoped().
		Delete(&entity.BlockEvent{}).Error

	if err != nil {
		tx.Rollback()
		return domain.NewError(500, err, nil)
	}

	circle.EventID = nil
	circle.Day = nil

	err = tx.Save(circle).Error

	if err != nil {
		tx.Rollback()
		return domain.NewError(500, err, nil)
	}

	tx.Commit()

	return nil
}

// OnboardNewCircle implements CircleRepo.
func (c *circleRepo) OnboardNewCircle(circle *entity.Circle, user *entity.User) (*entity.Circle, *domain.Error) {
	tx := c.db.Begin()
	if tx.Error != nil {
		return nil, domain.NewError(500, tx.Error, nil)
	}

	err := tx.Create(circle).Error
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(500, err, nil)
	}

	if circle.ID == 0 {
		tx.Rollback()
		return nil, domain.NewError(500, errors.New("CIRCLE_ID_NOT_FOUND"), nil)
	}

	user.CircleID = &circle.ID

	err = tx.Save(user).Error
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(500, err, nil)
	}

	tx.Commit()

	return circle, nil
}

// transformBlockStringIntoBlockEvent implements CircleRepo.
func (c *circleRepo) transformBlockStringIntoBlockEvent(block string) (*entity.BlockEvent, *domain.Error) {

	var postfix string
	var prefix string
	var name string

	splitted := strings.SplitN(block, "-", 2)
	if len(splitted) != 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	prefix = strings.ToUpper(splitted[0])
	postfix = strings.ToLower(splitted[1])
	name = strings.ToUpper(splitted[0]) + "-" + strings.ToLower(splitted[1])

	// check prefix length
	if len(postfix) > 10 || len(prefix) > 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	return &entity.BlockEvent{
		Prefix:  prefix,
		Postfix: postfix,
		Name:    name,
	}, nil
}

// UpserstOneCircle implements CircleRepo.
func (c *circleRepo) UpserstOneCircle(circle *entity.Circle) (*entity.Circle, *domain.Error) {
	err := c.db.Save(circle).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circle, nil
}

// UpdateCircleAndAllRelation implements CircleRepo.
func (c *circleRepo) UpdateCircleAndAllRelation(userID int, payload *entity.Circle, body *circle_dto.UpdateCircleRequestBody) ([]entity.CircleRaw, *domain.Error) {
	tx := c.db.Begin()
	if tx.Error != nil {
		return nil, domain.NewError(500, tx.Error, nil)
	}

	if body.FandomIDs != nil {
		if len(*body.FandomIDs) == 0 {
			// delete all fandom by circle id
			err := tx.Where("circle_id = ?", payload.ID).Delete(&entity.CircleFandom{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else if len(*body.FandomIDs) > 5 {
			tx.Rollback()
			return nil, domain.NewError(400, errors.New("FANDOM_LIMIT_EXCEEDED"), nil)
		} else {
			// delete all fandom by circle id
			err := tx.Where("circle_id = ?", payload.ID).Delete(&entity.CircleFandom{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			var circleFandoms []entity.CircleFandom
			for _, fandomID := range *body.FandomIDs {
				circleFandoms = append(circleFandoms, entity.CircleFandom{
					CircleID: payload.ID,
					FandomID: fandomID,
				})
			}

			err = tx.Create(&circleFandoms).Error

			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		}

	}

	if body.WorkTypeIDs != nil {
		if len(*body.WorkTypeIDs) == 0 {
			// delete all work type by circle id
			err := tx.Where("circle_id = ?", payload.ID).Delete(&entity.CircleWorkType{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else if len(*body.WorkTypeIDs) > 5 {
			tx.Rollback()
			return nil, domain.NewError(400, errors.New("WORK_TYPE_LIMIT_EXCEEDED"), nil)
		} else {

			// delete all work type by circle id
			err := tx.Where("circle_id = ?", payload.ID).Delete(&entity.CircleWorkType{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			var circleWorkTypes []entity.CircleWorkType
			for _, workTypeID := range *body.WorkTypeIDs {
				circleWorkTypes = append(circleWorkTypes, entity.CircleWorkType{
					CircleID:   payload.ID,
					WorkTypeID: workTypeID,
				})
			}

			err = tx.Create(&circleWorkTypes).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

		}
	}

	saveErr := tx.Save(&payload).Error
	if saveErr != nil {
		tx.Rollback()
		return nil, domain.NewError(500, saveErr, nil)
	}

	tx.Commit()

	rows, err := c.FindOneBySlugAndRelatedTables(payload.Slug, userID)
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(err.Code, err.Err, nil)
	}

	return rows, nil
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
	var row []entity.CircleRaw
	err := c.db.Select(`
			c.*,
			f.id as fandom_id,
			f.name as fandom_name,
			f.visible as fandom_visible,
			f.created_at as fandom_created_at,
			f.updated_at as fandom_updated_at,
			f.deleted_at as fandom_deleted_at,

			wt.id as work_type_id,
			wt.name as work_type_name,
			wt.created_at as work_type_created_at,
			wt.updated_at as work_type_updated_at,
			wt.deleted_at as work_type_deleted_at,

			p.id as product_id,
			p.name as product_name,
			p.image_url as product_image_url,
			p.created_at as product_created_at,
			p.updated_at as product_updated_at,
			p.deleted_at as product_deleted_at,

			e.name as event_name,
			e.slug as event_slug,
			e.description as event_description,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at,

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			user_bookmark.created_at as bookmarked_at,
			CASE WHEN user_bookmark.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS bookmarked
 		 `).
		Table("circle c").
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON cf.fandom_id = f.id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON cwt.work_type_id = wt.id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN user_bookmark ON c.id = user_bookmark.circle_id AND user_bookmark.user_id = COALESCE(?, user_bookmark.user_id)", userID).
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id").
		Where("c.deleted_at is null AND c.slug = ?", slug).Find(&row).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return row, nil
}

// findAllBookmarkedCount implements CircleRepo.
func (c *circleRepo) FindAllBookmarkedCount(userID int, filter *circle_dto.FindAllCircleFilter) (int, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	var count int64

	err := c.db.
		Select("count(DISTINCT c.id)").
		Table("circle c").
		Joins("JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", userID).
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id").
		Where(whereClause, args...).
		Count(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// FindBookmarkedCircleByUserID implements CircleRepo.
func (c *circleRepo) FindBookmarkedCircleByUserID(userID int, filter *circle_dto.FindAllCircleFilter) ([]entity.CircleRaw, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	cte := c.db.
		Select(`
			c.*,
			ub.created_at as bookmarked_at,
			true as bookmarked
		`).
		Table("circle c").
		Joins("JOIN user_bookmark ub on c.id = ub.circle_id AND ub.user_id = ?", userID).
		Where("c.deleted_at is null").
		Order("ub.created_at desc").
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit)

	join := c.db.Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	var circleRaw []entity.CircleRaw
	err := join.Table("(?) as c", cte).
		Select(`
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

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			e.name as event_name,
			e.slug as event_slug,
			e.description as event_description,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at
		`).
		Where(whereClause, args...).
		Find(&circleRaw).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circleRaw, nil
}

// findAllWhereSQL implements CircleRepo.
func (c *circleRepo) findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{}) {
	whereClause := `1 = 1`
	args := make([]interface{}, 0)

	if filter.Search != "" {
		whereClause += " and (c.name ILIKE ? OR f.name ILIKE ? OR wt.name ILIKE ? OR be.name ILIKE ?)"
		searchClause := fmt.Sprintf("%%%s%%", filter.Search)
		args = append(args, searchClause, searchClause, searchClause, searchClause)
	}

	if len(filter.FandomIDs) > 0 {
		whereClause += " and f.id in (?)"
		args = append(args, filter.FandomIDs)
	}

	if len(filter.WorkTypeIDs) > 0 {
		whereClause += " and wt.id in (?)"
		args = append(args, filter.WorkTypeIDs)
	}

	return whereClause, args
}

// FindAll implements CircleRepo.
func (c *circleRepo) FindAllCircles(filter *circle_dto.FindAllCircleFilter, userID int) ([]entity.CircleRaw, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	cte := c.db.Table("circle").Where("deleted_at IS NULL").Limit(filter.Limit).Offset((filter.Page - 1) * filter.Limit)

	var circles []entity.CircleRaw
	joins := c.db.
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", userID).
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	db := joins.
		Table("(?) as c", cte).
		Select(`
			c.*,
			f.id as fandom_id,
			f.name as fandom_name,
			f.visible as fandom_visible,
			f.created_at as fandom_created_at,
			f.updated_at as fandom_updated_at,
			f.deleted_at as fandom_deleted_at,

			wt.id as work_type_id,
			wt.name as work_type_name,
			wt.created_at as work_type_created_at,
			wt.updated_at as work_type_updated_at,
			wt.deleted_at as work_type_deleted_at,

			p.id as product_id,
			p.name as product_name,
			p.image_url as product_image_url,
			p.created_at as product_created_at,
			p.updated_at as product_updated_at,
			p.deleted_at as product_deleted_at,

			e.name as event_name,
			e.slug as event_slug,
			e.description as event_description,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at,

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			ub.created_at as bookmarked_at,
			CASE WHEN ub.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS bookmarked
		`).
		Where(whereClause, args...).
		Order("c.created_at desc")

	err := db.Find(&circles).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circles, nil
}

// FindAllCount implements CircleRepo.
func (c *circleRepo) FindAllCount(filter *circle_dto.FindAllCircleFilter) (int, *domain.Error) {
	whereClause, args := c.findAllWhereSQL(filter)

	var count int64
	err := c.db.
		Select("count(DISTINCT c.id)").
		Table("circle c").
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", 0).
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id").
		Where(whereClause, args...).
		Count(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// CreateOne implements CircleRepo.
func (c *circleRepo) CreateOne(circle entity.Circle) (*entity.Circle, *domain.Error) {
	err := c.db.Create(&circle).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &circle, nil
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
