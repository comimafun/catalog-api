package event_dto

type GetPaginatedEventsFilter struct {
	Page  int `query:"page" validate:"required,min=1"`
	Limit int `query:"limit" validate:"required,min=1,max=20"`
}

type CreateEventReqeuestBody struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	StartedAt   string  `json:"started_at" validate:"required,datetime=2006-01-02T15:04:05Z"`
	EndedAt     string  `json:"ended_at" validate:"required,datetime=2006-01-02T15:04:05Z"`
	Description *string `json:"description" validate:"omitempty"`
}
