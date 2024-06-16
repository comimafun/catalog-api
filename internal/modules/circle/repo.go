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
	FindOneByID(id int) (*entity.Circle, *domain.Error)
	FindOneBySlugAndRelatedTables(slug string, userID int) ([]entity.CircleRaw, *domain.Error)

	UpserstOneCircle(circle *entity.Circle) (*entity.Circle, *domain.Error)
	UpdateCircleAndAllRelation(circleID int, payload *entity.Circle, body *circle_dto.UpdateCircleRequestBody) ([]entity.CircleRaw, *domain.Error)

	transformBlockStringIntoCircleBlock(block string) (*entity.CircleBlock, *domain.Error)
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

// transformBlockStringIntoCircleBlock implements CircleRepo.
func (c *circleRepo) transformBlockStringIntoCircleBlock(block string) (*entity.CircleBlock, *domain.Error) {

	var postfix string
	var prefix string

	splitted := strings.SplitN(block, "-", 2)
	if len(splitted) != 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	prefix = strings.ToUpper(splitted[0])
	postfix = strings.ToLower(splitted[1])

	// check prefix length
	if len(postfix) > 8 || len(prefix) > 2 {
		return nil, domain.NewError(400, errors.New("INVALID_BLOCK_FORMAT"), nil)
	}

	return &entity.CircleBlock{
		Prefix:  prefix,
		Postfix: postfix,
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
func (c *circleRepo) UpdateCircleAndAllRelation(circleID int, payload *entity.Circle, body *circle_dto.UpdateCircleRequestBody) ([]entity.CircleRaw, *domain.Error) {
	tx := c.db.Begin()
	if tx.Error != nil {
		return nil, domain.NewError(500, tx.Error, nil)
	}

	if body.CircleBlock != nil && payload.EventID != nil {
		trimmedBlockString := strings.TrimSpace(*body.CircleBlock)
		if trimmedBlockString == "" {
			err := tx.Table("block_event").Where("circle_id = ?", circleID).Unscoped().Delete(&entity.BlockEvent{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else {
			block, err := c.transformBlockStringIntoBlockEvent(trimmedBlockString)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			existingBlock := new(entity.BlockEvent)
			existingErr := tx.Where("prefix = ? AND postfix = ? AND circle_id = ? AND event_id = ?", block.Prefix, block.Postfix, circleID, payload.EventID).First(&existingBlock).Error
			if existingErr != nil && !errors.Is(existingErr, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, domain.NewError(500, existingErr, nil)
			}

			if existingBlock.ID != 0 {
				tx.Rollback()
				return nil, domain.NewError(400, errors.New("BLOCK_ALREADY_EXIST"), nil)
			}

			// delete block
			deleteErr := tx.Table("block_event").Where("circle_id = ?", circleID).Unscoped().Delete(&entity.BlockEvent{}).Error
			if deleteErr != nil {
				tx.Rollback()
				return nil, domain.NewError(500, deleteErr, nil)
			}

			block.CircleID = circleID
			block.EventID = *payload.EventID

			createErr := tx.Create(block).Error
			if createErr != nil {
				return nil, domain.NewError(500, createErr, nil)
			}
		}

	}

	if body.FandomIDs != nil {
		if len(*body.FandomIDs) == 0 {
			// delete all fandom by circle id
			err := tx.Where("circle_id = ?", circleID).Delete(&entity.CircleFandom{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else if len(*body.FandomIDs) > 5 {
			tx.Rollback()
			return nil, domain.NewError(400, errors.New("FANDOM_LIMIT_EXCEEDED"), nil)
		} else {
			// delete all fandom by circle id
			err := tx.Where("circle_id = ?", circleID).Delete(&entity.CircleFandom{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			var circleFandomIDs []string
			var circleFandomArgs []interface{}
			for _, fandomID := range *body.FandomIDs {
				circleFandomIDs = append(circleFandomIDs, "(?,?)")
				circleFandomArgs = append(circleFandomArgs, circleID, fandomID)
			}

			query := fmt.Sprintf(
				`
	INSERT INTO circle_fandom (circle_id, fandom_id)
	VALUES %s
`, strings.Join(circleFandomIDs, ", "))

			err = tx.Exec(query, circleFandomArgs...).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		}

	}

	if body.WorkTypeIDs != nil {
		if len(*body.WorkTypeIDs) == 0 {
			// delete all work type by circle id
			err := tx.Where("circle_id = ?", circleID).Delete(&entity.CircleWorkType{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else if len(*body.WorkTypeIDs) > 5 {
			tx.Rollback()
			return nil, domain.NewError(400, errors.New("WORK_TYPE_LIMIT_EXCEEDED"), nil)
		} else {

			// delete all work type by circle id
			err := tx.Where("circle_id = ?", circleID).Delete(&entity.CircleWorkType{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			var circleWorkTypeIDs []string
			var circleWorkTypeArgs []interface{}
			for _, workTypeID := range *body.WorkTypeIDs {
				circleWorkTypeIDs = append(circleWorkTypeIDs, "(?,?)")
				circleWorkTypeArgs = append(circleWorkTypeArgs, circleID, workTypeID)
			}

			query := fmt.Sprintf(`
			INSERT INTO circle_work_type (circle_id, work_type_id)
			VALUES %s
		`, strings.Join(circleWorkTypeIDs, ", "))

			err = tx.Exec(query, circleWorkTypeArgs...).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		}
	}

	if body.Products != nil {
		if len(*body.Products) == 0 {
			// delete all products by circle id
			err := tx.Where("circle_id = ?", circleID).Delete(&entity.Product{}).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		} else if len(*body.Products) > 5 {
			tx.Rollback()
			return nil, domain.NewError(400, errors.New("PRODUCT_LIMIT_EXCEEDED"), nil)
		} else {
			inputs := make([]entity.Product, 0)
			for _, product := range *body.Products {
				inputs = append(inputs, entity.Product{
					ID:       product.ID,
					Name:     product.Name,
					ImageURL: product.ImageURL,
					CircleID: circleID,
				})

			}

			var existingProducts []entity.Product
			err := tx.Where("circle_id = ?", circleID).Find(&existingProducts).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			newProductMap := make(map[int]bool)

			for _, input := range inputs {
				if input.ID == 0 {
					err := tx.Create(&input).Error
					if err != nil {
						tx.Rollback()
						return nil, domain.NewError(500, err, nil)
					}

					newProductMap[input.ID] = true
				} else {
					err := tx.Save(&input).Error
					if err != nil {
						tx.Rollback()
						return nil, domain.NewError(500, err, nil)
					}

					newProductMap[input.ID] = true
				}
			}

			for _, product := range existingProducts {
				_, ok := newProductMap[product.ID]
				if !ok {
					err := tx.Delete(&entity.Product{}, product.ID).Error
					if err != nil {
						tx.Rollback()
						return nil, domain.NewError(500, err, nil)
					}
				}
			}
		}

	}

	saveErr := tx.Save(&payload).Error
	if saveErr != nil {
		tx.Rollback()
		return nil, domain.NewError(500, saveErr, nil)
	}

	tx.Commit()

	rows, err := c.FindOneBySlugAndRelatedTables(payload.Slug, 0)
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

			ub.created_at AS bookmarked_at,
			CASE WHEN ub.user_id IS NOT NULL THEN
				TRUE
			ELSE
				FALSE
			END AS bookmarked,

			e."name" as event_name,
			e.slug as event_slug,
			e.description as event_description,
			e.started_at as event_started_at,
			e.ended_at as event_ended_at,

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name
		FROM
			circle c
		LEFT JOIN circle_fandom cf ON c.id = cf.circle_id
		LEFT JOIN fandom f ON f.id = cf.fandom_id
		LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id
		LEFT JOIN work_type wt ON wt.id = cwt.work_type_id
		LEFT JOIN product p on c.id = p.circle_id
		LEFT JOIN user_bookmark ub ON c.id = ub.circle_id AND ub.user_id = COALESCE(?, ub.user_id)
		LEFT JOIN "event" e ON c.event_id = e.id
		LEFT JOIN block_event be on c.id = be.circle_id
		WHERE
				c.deleted_at IS NULL AND
				c.slug = ?
		`

	var row []entity.CircleRaw

	err := c.db.Raw(query, userID, slug).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, err, nil)
		}
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
		WITH PaginatedCircle AS (
			SELECT
				c.*
			FROM
				circle c
			WHERE
				c.deleted_at IS NULL
			LIMIT ?
			OFFSET ?
		)

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

			be.id as block_event_id,
			be.prefix as block_event_prefix,
			be.postfix as block_event_postfix,
			be.name as block_event_name,

			ub.created_at as bookmarked_at,
			CASE WHEN ub.user_id is not null THEN true ELSE false END as bookmarked
		from 
			PaginatedCircle c
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
			block_event be on c.id = be.circle_id
		%s
		ORDER by
			ub.created_at desc
	`, whereClause)

	args = append(args, userID)

	offset := (filter.Page - 1) * filter.Limit
	args = append([]interface{}{filter.Limit, offset, userID}, args...)

	var circles []entity.CircleRaw
	err := c.db.Raw(query, args...).Scan(&circles).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return circles, nil
}

// findAllWhereSQL implements CircleRepo.
func (c *circleRepo) findAllWhereSQL(filter *circle_dto.FindAllCircleFilter) (string, []interface{}) {
	whereClause := `WHERE
						1 = 1`
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

	query := fmt.Sprintf(`
		WITH PaginatedCircle AS (
			SELECT
				c.*
			FROM
				circle c
			WHERE
				c.deleted_at IS NULL
			LIMIT ? OFFSET ?
		)
		SELECT
			c.*,
			f.name AS fandom_name,
			f.id AS fandom_id,
			f.visible AS fandom_visible,
			f.created_at AS fandom_created_at,
			f.updated_at AS fandom_updated_at,
			f.deleted_at AS fandom_deleted_at,
			wt.name AS work_type_name,
			wt.id AS work_type_id,
			wt.created_at AS work_type_created_at,
			wt.updated_at AS work_type_updated_at,
			wt.deleted_at AS work_type_deleted_at,

			be.id AS block_event_id,
			be.prefix AS block_event_prefix,
			be.postfix AS block_event_postfix,
			be.name AS block_event_name,
			
			ub.created_at AS bookmarked_at,
			CASE WHEN ub.user_id IS NOT NULL THEN
				TRUE
			ELSE
				FALSE
			END AS bookmarked
		FROM
			PaginatedCircle c
			LEFT JOIN circle_fandom cf ON c.id = cf.circle_id
			LEFT JOIN fandom f ON f.id = cf.fandom_id
			LEFT JOIN circle_work_type cwt ON c.id = cwt.circle_id
			LEFT JOIN work_type wt ON wt.id = cwt.work_type_id
			LEFT JOIN user_bookmark ub ON c.id = ub.circle_id
				AND ub.user_id = COALESCE(?, ub.user_id)
			LEFT JOIN block_event be ON c.id = be.circle_id
		%s
		ORDER BY
				c.created_at desc
	`, whereClause)

	offset := (filter.Page - 1) * filter.Limit
	args = append([]interface{}{filter.Limit, offset, userID}, args...)

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

// FindOneByID implements CircleRepo.
func (c *circleRepo) FindOneByID(id int) (*entity.Circle, *domain.Error) {
	var circle entity.Circle
	err := c.db.First(&circle, id).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &circle, nil
}
