package report

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type ReportService struct {
	repo *ReportRepo
}

// Initialize Circle Report Service
func NewReportService(repo *ReportRepo) *ReportService {
	return &ReportService{repo}
}

// CreateReport implements CircleReportService
func (r *ReportService) CreateReportCircle(report *entity.Report) *domain.Error {
	return r.repo.CreateReportCircle(report)
}

// FindCircleReportByID implements CircleReportService
func (r *ReportService) FindReportByID(id int) (*entity.Report, *domain.Error) {
	return r.repo.FindByID(id)
}

// FindAllCircleReportByCircleID implements CircleReportService
func (r *ReportService) FindAllReportByCircleID(circleID int) ([]entity.Report, *domain.Error) {
	return r.repo.FindAllByCircleID(circleID)
}
