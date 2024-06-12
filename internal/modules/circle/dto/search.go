package circle_dto

type FindAllCircleFilter struct {
	Search      string `query:"search" validate:"omitempty"`
	WorkTypeIDs []int  `query:"work_type_id" validate:"omitempty,dive"`
	Page        int    `query:"page" validate:"required,min=1"`
	Limit       int    `query:"limit" validate:"required,min=1,max=20"`
	FandomIDs   []int  `query:"fandom_id" validate:"omitempty,dive"`
}
