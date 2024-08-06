package circle_report_dto

type CreateCircleReportPayload struct {
	Reason string `json:"reason" validate:"required,min=3,max=255"`
}
