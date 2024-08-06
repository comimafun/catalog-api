package circle_report_dto

type CreateCircleReportBody struct {
	Reason string `json:"reason" validate:"required,min=3"`
}
