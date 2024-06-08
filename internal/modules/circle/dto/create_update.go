package circle_dto

import (
	"catalog-be/internal/entity"
	"time"
)

type ImageURLs struct {
	URL          string `json:"url" validate:"omitempty,url,max=255"`
	PictureURL   string `json:"picture_url" validate:"omitempty,url,max=255"`
	FacebookURL  string `json:"facebook_url" validate:"omitempty,url,max=255"`
	InstagramURL string `json:"instagram_url" validate:"omitempty,url,max=255"`
	TwitterURL   string `json:"twitter_url" validate:"omitempty,url,max=255"`
}

type OnboardNewCircleRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	ImageURLs
}

type CreateFandomCircleRelation struct {
	ID   int    `json:"ID" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type CreateWorkTypeCircleRelation struct {
	ID   int    `json:"ID" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type UpdateCircleRequestBody struct {
	// Name        *string `json:"name" validate:"omitempty,min=3,max=255"`
	CircleBlock string      `json:"circle_block" validate:"omitempty"`
	Description string      `json:"description" validate:"omitempty"`
	Batch       *int        `json:"batch" validate:"omitempty"`
	Day         *entity.Day `json:"day" validate:"omitempty,oneof=first second both"`
	ImageURLs
	FandomIDs   []int `json:"fandom_ids" validate:"omitempty,dive"`
	WorkTypeIDs []int `json:"work_type_ids" validate:"omitempty,dive"`
}

type CircleResponse struct {
	entity.Circle
	Fandom   []entity.Fandom   `json:"fandom"`
	WorkType []entity.WorkType `json:"work_type"`

	Bookmarked   bool       `json:"bookmarked"`
	BookmarkedAt *time.Time `json:"bookmarked_at"`
}
