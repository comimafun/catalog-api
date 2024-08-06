package report_dto

type CreateReportPayload struct {
	Reason string `json:"reason" validate:"required,min=3,max=255"`
}
