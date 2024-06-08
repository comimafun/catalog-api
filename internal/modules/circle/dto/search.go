package circle_dto

type FindAllCircleFilter struct {
	Search     string `json:"search" validate:"omitempty"`
	FandomID   int    `json:"fandom_id" validate:"omitempty,min=1"`
	WorkTypeID int    `json:"work_type_id" validate:"omitempty,min=1"`
	Page       int    `json:"page" validate:"required,min=1"`
	Limit      int    `json:"limit" validate:"required,min=1,max=20"`
}
