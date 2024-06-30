package circle_dto

import (
	"catalog-be/internal/entity"
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

type ProductBody struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=255"`
	ImageURL string `json:"image_url" validate:"required,url"`
}

type UpdateCircleRequestBody struct {
	Name        *string     `json:"name" validate:"omitempty,min=3,max=255"`
	CircleBlock *string     `json:"circle_block" validate:"omitempty"`
	Description *string     `json:"description" validate:"omitempty"`
	Batch       *int        `json:"batch" validate:"omitempty"`
	Day         *entity.Day `json:"day" validate:"omitempty,day_or_empty"`

	PictureURL      *string `json:"picture_url" validate:"omitempty,max=255"`
	CoverPictureURL *string `json:"cover_picture_url" validate:"omitempty,max=255"`

	URL          *string `json:"url" validate:"omitempty,url_or_empty,max=255"`
	FacebookURL  *string `json:"facebook_url" validate:"omitempty,url_or_empty,max=255"`
	InstagramURL *string `json:"instagram_url" validate:"omitempty,url_or_empty,max=255"`
	TwitterURL   *string `json:"twitter_url" validate:"omitempty,url_or_empty,max=255"`

	EventID     *int           `json:"event_id" validate:"omitempty"`
	FandomIDs   *[]int         `json:"fandom_ids" validate:"omitempty,dive"`
	WorkTypeIDs *[]int         `json:"work_type_ids" validate:"omitempty,dive"`
	Products    *[]ProductBody `json:"products" validate:"omitempty,dive"`
}

type BlockResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CircleResponse struct {
	entity.Circle
	Fandom   []entity.Fandom   `json:"fandom"`
	WorkType []entity.WorkType `json:"work_type"`
	Product  []entity.Product  `json:"product"`

	Bookmarked bool `json:"bookmarked"`

	BlockEvent *BlockResponse `json:"block"`

	Event *entity.Event `json:"event"`
}

type CircleOneForPaginationResponse struct {
	entity.Circle
	Fandom   []entity.Fandom   `json:"fandom"`
	WorkType []entity.WorkType `json:"work_type"`
	Product  []entity.Product  `json:"product"`

	Bookmarked bool           `json:"bookmarked"`
	BlockEvent *BlockResponse `json:"block"`
}
