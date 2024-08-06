package report

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/circle"
	"errors"

	"gorm.io/gorm"
)

type ReportService struct {
	repo       *ReportRepo
	circleRepo *circle.CircleRepo
}

// Initialize Circle Report Service
func NewReportService(repo *ReportRepo, circleRepo *circle.CircleRepo) *ReportService {
	return &ReportService{
		repo,
		circleRepo,
	}
}

// CreateReport implements CircleReportService
func (r *ReportService) CreateReportCircle(report *entity.Report) *domain.Error {
	if _, err := r.circleRepo.GetOneCircleByCircleID(report.CircleID); err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return domain.NewError(404, errors.New("CIRCLE_NOT_FOUND"), nil)
		}
		return err
	}
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
