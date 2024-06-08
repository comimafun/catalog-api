package circle_work_type

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CircleWorkTypeRepo interface {
	deleteWorkTypeRelationByCircleID(circleID int) *domain.Error
	BatchInsertCircleWorkTypeRelation(circleID int, workTypeIDs []int) *domain.Error
	FindAllCircleRelationWorkType(circleID int) ([]entity.WorkType, *domain.Error)
}

type circleWorkTypeRepo struct {
	db *gorm.DB
}

// BatchInsertCircleWorkTypeRelation implements CircleWorkTypeRepo.
func (c *circleWorkTypeRepo) BatchInsertCircleWorkTypeRelation(circleID int, workTypeIDs []int) *domain.Error {
	tx := c.db.Begin()
	if tx.Error != nil {
		return domain.NewError(500, tx.Error, nil)
	}

	deleteErr := c.deleteWorkTypeRelationByCircleID(circleID)

	if deleteErr != nil {
		tx.Rollback()
		return deleteErr
	}

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

	err := tx.Exec(query, valueArgs...).Error
	if err != nil {
		tx.Rollback()
		return domain.NewError(500, err, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return domain.NewError(500, err, nil)
	}

	return nil
}

// FindAllCircleRelationWorkType implements CircleWorkTypeRepo.
func (c *circleWorkTypeRepo) FindAllCircleRelationWorkType(circleID int) ([]entity.WorkType, *domain.Error) {
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

// deleteWorkTypeRelationByCircleID implements CircleWorkTypeRepo.
func (c *circleWorkTypeRepo) deleteWorkTypeRelationByCircleID(circleID int) *domain.Error {
	err := c.db.Exec(`
    delete from circle_work_type where circle_id = ?
    `, circleID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

func NewCircleWorkTypeRepo(db *gorm.DB) CircleWorkTypeRepo {
	return &circleWorkTypeRepo{db: db}
}
