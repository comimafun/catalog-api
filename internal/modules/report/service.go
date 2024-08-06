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
func (r *ReportService) CreateCircleReport(circleID int, userID int, reason string) *domain.Error {
	return r.repo.CreateReport(circleID, userID, reason)
}

// FindCircleReportByID implements CircleReportService
func (r *ReportService) FindCircleReportByID(id int) (*entity.CircleReport, *domain.Error) {
	return r.repo.FindByID(id)
}

// FindAllCircleReportByCircleID implements CircleReportService
func (r *ReportService) FindAllCircleReportByCircleID(circleID int) ([]entity.CircleReport, *domain.Error) {
	return r.repo.FindByCircleIDAndUserID(circleID, 0)
}
