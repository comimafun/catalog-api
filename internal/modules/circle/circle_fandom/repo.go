package circle_fandom

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CircleFandomRepo struct {
	db *gorm.DB
}

// BatchInsertFandomCircleRelation implements CircleFandomRepo.
func (c *CircleFandomRepo) BatchInsertFandomCircleRelation(circleID int, fandomIDs []int) *domain.Error {
	tx := c.db.Begin()
	if tx.Error != nil {
		return domain.NewError(500, tx.Error, nil)
	}

	deleteErr := c.deleteFandomRelationByCircleID(circleID)
	if deleteErr != nil {
		tx.Rollback()
		return deleteErr
	}
	var valueStrings []string
	valueArgs := make([]interface{}, 0)
	for _, fandomID := range fandomIDs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, circleID, fandomID)
	}

	query := fmt.Sprintf(`
			INSERT INTO circle_fandom (circle_id, fandom_id)
			VALUES %s
		`, strings.Join(valueStrings, ", "))

	err := c.db.Exec(query, valueArgs...).Error
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

// FindAllCircleRelationFandom implements CircleFandomRepo.
func (c *CircleFandomRepo) FindAllCircleRelationFandom(circleID int) ([]entity.Fandom, *domain.Error) {
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

// deleteFandomRelationByCircleID implements CircleFandomRepo.
func (c *CircleFandomRepo) deleteFandomRelationByCircleID(circleID int) *domain.Error {
	err := c.db.Exec(`
    delete from circle_fandom where circle_id = ?
`, circleID).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

func NewCircleFandomRepo(db *gorm.DB) *CircleFandomRepo {
	return &CircleFandomRepo{
		db: db,
	}
}
