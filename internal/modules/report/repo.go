package report

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ReportRepo struct {
	db *gorm.DB
}

// Initialize
func NewReportRepo(db *gorm.DB) *ReportRepo {
	return &ReportRepo{db}
}

// Create Report for Circle
func (r *ReportRepo) CreateReport(circleID int, userID int, reason string) *domain.Error {
	err := r.db.Table("circle_report").Create(&entity.CircleReport{
		UserID:   userID,
		CircleID: circleID,
		Reason:   reason,
	}).Error

	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// Find certain Report by ID
func (r *ReportRepo) FindByID(id int) (*entity.CircleReport, *domain.Error) {
	report := new(entity.CircleReport)
	err := r.db.Table("circle_report").Where("id = ?", id).First(report).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return report, nil
}

// Find All Report by Circle ID and User ID
func (r *ReportRepo) FindByCircleIDAndUserID(circleID int, userID int) ([]entity.CircleReport, *domain.Error) {
	var reports []entity.CircleReport
	err := r.db.Table("circle_report").Where("circle_id = ? AND user_id = ?", circleID, userID).Find(&reports).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return reports, nil
}
