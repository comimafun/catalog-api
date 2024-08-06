package circle_dto

import "catalog-be/internal/entity"

type GetPaginatedCirclesFilter struct {
	Search      string      `query:"search" validate:"omitempty"`
	WorkTypeIDs []int       `query:"work_type_id" validate:"omitempty,dive"`
	Page        int         `query:"page" validate:"required,min=1"`
	Limit       int         `query:"limit" validate:"required,min=1,max=20"`
	FandomIDs   []int       `query:"fandom_id" validate:"omitempty,dive"`
	Rating      []string    `query:"rating" validate:"omitempty,dive,oneof=GA PG M"`
	Event       string      `query:"event" validate:"omitempty"`
	Day         *entity.Day `query:"day" validate:"omitempty,oneof=first second both"`
}
