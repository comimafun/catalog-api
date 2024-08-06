package circle_report

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type CircleReportService struct {
	repo *CircleReportRepo
}

// Initialize Circle Report Service
func NewCircleReportService(repo *CircleReportRepo) *CircleReportService {
	return &CircleReportService{repo}
}

// CreateReport implements CircleReportService
func (c *CircleReportService) CreateCircleReport(circleID int, userID int, reason string) *domain.Error {
	return c.repo.CreateReport(circleID, userID, reason)
}

// FindCircleReportByID implements CircleReportService
func (c *CircleReportService) FindCircleReportByID(id int) (*entity.CircleReport, *domain.Error) {
	return c.repo.FindByID(id)
}

// FindAllCircleReportByCircleID implements CircleReportService
func (c *CircleReportService) FindAllCircleReportByCircleID(circleID int) ([]entity.CircleReport, *domain.Error) {
	return c.repo.FindByCircleIDAndUserID(circleID, 0)
}
