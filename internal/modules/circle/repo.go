package circle

import (
	"catalog-be/internal/domain"
	"errors"
	"os"
	"strings"

	"catalog-be/internal/entity"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"fmt"

	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"gorm.io/gorm"
)

type CircleRepo struct {
	db *gorm.DB
}

// UpdateAttendingEventDayAndCircleBlock implements CircleRepo.
func (c *CircleRepo) UpdateAttendingEventDayAndCircleBlock(circle *entity.Circle, body *circle_dto.UpdateCircleAttendingEventDayAndBlockPayload) *domain.Error {
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

// DeleteAllBlockEventByCircleIDAndEventID implements CircleRepo.
func (c *CircleRepo) DeleteAllBlockEventByCircleIDAndEventID(circle *entity.Circle) *domain.Error {
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
func (c *CircleRepo) OnboardNewCircle(circle *entity.Circle, user *entity.User) (*entity.Circle, *domain.Error) {
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
func (c *CircleRepo) transformBlockStringIntoBlockEvent(block string) (*entity.BlockEvent, *domain.Error) {

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

// UpsertOneCircle implements CircleRepo.
func (c *CircleRepo) UpsertOneCircle(circle *entity.Circle) (*entity.Circle, *domain.Error) {
	err := c.db.Save(circle).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circle, nil
}

// UpdateOneCircleAndAllRelation implements CircleRepo.
func (c *CircleRepo) UpdateOneCircleAndAllRelation(userID int, payload *entity.Circle, body *circle_dto.UpdateCirclePayload) ([]entity.CircleJoinedTables, *domain.Error) {
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

	rows, err := c.GetOneCircleJoinTablesByCircleSlug(payload.Slug, userID)
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(err.Code, err.Err, nil)
	}

	return rows, nil
}

func NewCircleRepo(
	db *gorm.DB,
) *CircleRepo {
	return &CircleRepo{
		db: db,
	}
}

// GetOneCircleJoinTablesByCircleSlug implements CircleRepo.
func (c *CircleRepo) GetOneCircleJoinTablesByCircleSlug(slug string, userID int) ([]entity.CircleJoinedTables, *domain.Error) {
	var row []entity.CircleJoinedTables
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
func (c *CircleRepo) GetAllBookmarkedCircleCount(userID int, filter *circle_dto.GetPaginatedCirclesFilter) (int, *domain.Error) {
	var count int64

	dbs := c.db.
		Select("count(DISTINCT c.id)").
		Table("circle c").
		Joins("JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", userID).
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	err := dbs.
		Count(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// GetPaginatedBookmarkedCirclesByUserID implements CircleRepo.
func (c *CircleRepo) GetPaginatedBookmarkedCirclesByUserID(userID int, filter *circle_dto.GetPaginatedCirclesFilter) ([]entity.CircleJoinedTables, *domain.Error) {

	cte := c.db.
		Select(`
			c.id as id,
			c.name as name,
			c.slug as slug,
			c.picture_url as picture_url,
			c.url as url,
			c.facebook_url as facebook_url,
			c.twitter_url as twitter_url,
			c.instagram_url as instagram_url,
			c.verified as verified,
			c.published as published,
			c.created_at as created_at,
			c.updated_at as updated_at,
			c.deleted_at as deleted_at,
			c.day as day,
			c.event_id as event_id,
			c.cover_picture_url as cover_picture_url,
			c.rating as rating,
			ub.created_at as bookmarked_at,
			true as bookmarked
		`).
		Table("circle c").
		Joins("JOIN user_bookmark ub on c.id = ub.circle_id AND ub.user_id = ?", userID).
		Where("c.deleted_at is null").
		Order("ub.created_at desc").
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit)

	join := c.db.Clauses(exclause.NewWith("cte", cte)).Table("cte as c")

	join = join.Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	var circleRaw []entity.CircleJoinedTables
	err := join.
		Select(`
			c.*,

			f."name" as fandom_name,
			f.id as fandom_id,
			f.visible as fandom_visible,

			wt."name" as work_type_name,
			wt.id as work_type_id,


			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			e.name as event_name,
			e.slug as event_slug,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at
		`).
		Order("c.bookmarked_at desc").
		Find(&circleRaw).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circleRaw, nil
}

// FindAll implements CircleRepo.
func (c *CircleRepo) GetPaginatedCircles(filter *circle_dto.GetPaginatedCirclesFilter, userID int) ([]entity.CircleJoinedTables, *domain.Error) {
	appStage := os.Getenv("APP_STAGE")

	cte := c.db.
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id")

	cte = cte.Table("circle c").
		Where("c.deleted_at IS NULL").
		Where("c.verified IS TRUE")

	if len(filter.Rating) > 0 {
		cte = cte.Where("c.rating IN (?)", filter.Rating)
	}
	if filter.Event != "" {
		cte = cte.Where("e.slug = ?", filter.Event)
	}

	if len(filter.FandomIDs) > 0 {
		cte = cte.Where("f.id in (?)", filter.FandomIDs)
	}

	if len(filter.WorkTypeIDs) > 0 {
		cte = cte.Where("wt.id in (?)", filter.WorkTypeIDs)
	}

	if filter.Day != nil {
		cte = cte.Where("c.day = ?", filter.Day)
	}

	if filter.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", filter.Search)
		cte = cte.Where("c.name ILIKE ? OR f.name ILIKE ? OR wt.name ILIKE ? OR be.name ILIKE ?",
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery)
	}

	if appStage == "production" {
		cte = cte.Where("c.published IS TRUE")
	}

	cte = cte.Select("c.id")
	cte = cte.
		Distinct("c.id").
		Order("c.id desc").
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit)

	var circles []entity.CircleJoinedTables
	joins := c.db.Clauses(exclause.NewWith("cte", cte)).Table("cte as cte")
	joins = joins.
		Joins("INNER JOIN circle c ON cte.id = c.id").
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", userID).
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	joins = joins.
		Select(`
			c.id as id,
			c.name as name,
			c.slug as slug,
			c.picture_url as picture_url,
			c.url as url,
			c.facebook_url as facebook_url,
			c.twitter_url as twitter_url,
			c.instagram_url as instagram_url,
			c.verified as verified,
			c.published as published,
			c.created_at as created_at,
			c.updated_at as updated_at,
			c.deleted_at as deleted_at,
			c.day as day,
			c.event_id as event_id,
			c.cover_picture_url as cover_picture_url,
			c.rating as rating,

			f.id as fandom_id,
			f.name as fandom_name,
			f.visible as fandom_visible,

			wt.id as work_type_id,
			wt.name as work_type_name,

			e.name as event_name,
			e.slug as event_slug,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at,

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			ub.created_at as bookmarked_at,
			CASE WHEN ub.user_id IS NOT NULL THEN TRUE ELSE FALSE END AS bookmarked
		`).
		Order("c.id desc")

	err := joins.Unscoped().Find(&circles).Error

	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circles, nil
}

// GetAllCirclesCount implements CircleRepo.
func (c *CircleRepo) GetAllCirclesCount(filter *circle_dto.GetPaginatedCirclesFilter) (int, *domain.Error) {
	appStage := os.Getenv("APP_STAGE")
	var count int64
	joins := c.db.
		Select("count(DISTINCT c.id)").
		Table("circle c").
		Joins("LEFT JOIN circle_fandom cf ON c.id = cf.circle_id").
		Joins("LEFT JOIN fandom f ON f.id = cf.fandom_id").
		Joins("LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id").
		Joins("LEFT JOIN work_type wt ON wt.id = cwt.work_type_id").
		Joins("LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)", 0).
		Joins("LEFT JOIN product p ON c.id = p.circle_id").
		Joins("LEFT JOIN event e ON c.event_id = e.id").
		Joins("LEFT JOIN block_event be ON c.id = be.circle_id AND be.event_id = c.event_id")

	if filter.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", filter.Search)
		joins = joins.Where("c.name ILIKE ? OR f.name ILIKE ? OR wt.name ILIKE ? OR be.name ILIKE ?",
			searchQuery,
			searchQuery,
			searchQuery,
			searchQuery)
	}

	if len(filter.Rating) > 0 {
		joins = joins.Where("c.rating IN (?)", filter.Rating)
	}
	if filter.Event != "" {
		joins = joins.Where("e.slug = ?", filter.Event)
	}

	if len(filter.FandomIDs) > 0 {
		joins = joins.Where("f.id in (?)", filter.FandomIDs)
	}

	if len(filter.WorkTypeIDs) > 0 {
		joins = joins.Where("wt.id in (?)", filter.WorkTypeIDs)
	}

	if filter.Day != nil {
		joins = joins.Where("c.day = ?", filter.Day)
	}

	if appStage == "production" {
		joins = joins.Where("c.published IS TRUE")
	}

	err := joins.
		Where("c.deleted_at is null and c.verified IS TRUE").
		Count(&count).Error

	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// GetOneCircleByCircleID implements CircleRepo.
func (c *CircleRepo) GetOneCircleByCircleID(id int) (*entity.Circle, *domain.Error) {
	var circle entity.Circle
	err := c.db.First(&circle, id).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &circle, nil
}
