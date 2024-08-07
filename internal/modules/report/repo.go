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
func (r *ReportRepo) CreateReportCircle(report *entity.Report) *domain.Error {
	err := r.db.Table("report").Create(report).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// Find certain Report by ID
func (r *ReportRepo) FindByID(id int) (*entity.Report, *domain.Error) {
	report := new(entity.Report)
	err := r.db.Table("report").Where("id = ?", id).First(report).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return report, nil
}

// Find All Report by Circle ID
func (r *ReportRepo) FindAllByCircleID(circleID int) ([]entity.Report, *domain.Error) {
	var reports []entity.Report
	err := r.db.Table("report").Where("circle_id = ?", circleID).Find(&reports).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return reports, nil
}
