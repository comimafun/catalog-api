package circle_dto

import "catalog-be/internal/entity"

type ImageURLs struct {
	PictureURL   *string `json:"picture_url" validate:"omitempty,url,max=255"`
	FacebookURL  *string `json:"facebook_url" validate:"omitempty,url,max=255"`
	InstagramURL *string `json:"instagram_url" validate:"omitempty,url,max=255"`
	TwitterURL   *string `json:"twitter_url" validate:"omitempty,url,max=255"`
}

type OnboardNewCircleRequestBody struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	ImageURLs
}

type UpdateCircleRequestBody struct {
	// Name        *string `json:"name" validate:"omitempty,min=3,max=255"`
	CircleBlock *string     `json:"circle_block" validate:"omitempty"`
	Description *string     `json:"description" validate:"omitempty"`
	Batch       *int        `json:"batch" validate:"omitempty"`
	Day         *entity.Day `json:"day" validate:"omitempty,oneof=first second both"`
	ImageURLs
}