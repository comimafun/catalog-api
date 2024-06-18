package entity

import (
	"time"

	"gorm.io/gorm"
)

// create enum day in go consisted of `first`, `second`, `both`
type Day string

const (
	First  Day = "first"
	Second Day = "second"
	Both   Day = "both"
)

type Circle struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	URL          string         `json:"url"`
	PictureURL   *string        `json:"picture_url"`
	FacebookURL  *string        `json:"facebook_url"`
	InstagramURL *string        `json:"instagram_url"`
	TwitterURL   *string        `json:"twitter_url"`
	Description  *string        `json:"description"`
	Batch        *int           `json:"batch"`
	Verified     bool           `json:"verified"`
	Published    bool           `json:"published"`
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`

	EventID *int `json:"event_id"`

	Day *Day `json:"day"`
}

type CircleRaw struct {
	Circle

	FandomID        int            `json:"fandom_id"`
	FandomName      string         `json:"fandom_name"`
	FandomVisible   bool           `json:"fandom_visible"`
	FandomCreatedAt *time.Time     `json:"fandom_created_at"`
	FandomUpdatedAt *time.Time     `json:"fandom_updated_at"`
	FandomDeletedAt gorm.DeletedAt `json:"-"`

	WorkTypeID        int            `json:"work_type_id"`
	WorkTypeName      string         `json:"work_type_name"`
	WorkTypeCreatedAt *time.Time     `json:"work_type_created_at"`
	WorkTypeUpdatedAt *time.Time     `json:"work_type_updated_at"`
	WorkTypeDeletedAt gorm.DeletedAt `json:"-"`

	Bookmarked   bool       `json:"bookmarked"`
	BookmarkedAt *time.Time `json:"bookmarked_at"`

	ProductID        int        `json:"product_id"`
	ProductName      string     `json:"product_name"`
	ProductImageURL  string     `json:"product_image_url"`
	ProductCreatedAt *time.Time `json:"product_created_at"`
	ProductUpdatedAt *time.Time `json:"product_updated_at"`

	EventName        string     `json:"event_name"`
	EventSlug        string     `json:"event_slug"`
	EventDescription string     `json:"event_description"`
	EventStartedAt   *time.Time `json:"event_started_at"`
	EventEndedAt     *time.Time `json:"event_ended_at"`

	BlockEventID      int    `json:"block_event_id"`
	BlockEventPrefix  string `json:"block_event_prefix"`
	BlockEventPostfix string `json:"block_event_postfix"`
	BlockEventName    string `json:"block_event_name"`
}

func (Circle) TableName() string {
	return "circle"
}

type CircleFandom struct {
	CircleID  int        `json:"circle_id"`
	FandomID  int        `json:"fandom_id"`
	CreatedAt *time.Time `json:"created_at"`
}

func (CircleFandom) TableName() string {
	return "circle_fandom"
}

type CircleWorkType struct {
	CircleID   int        `json:"circle_id"`
	WorkTypeID int        `json:"work_type_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

func (CircleWorkType) TableName() string {
	return "circle_work_type"
}

type UserBookmark struct {
	UserID    int        `json:"user_id"`
	CircleID  int        `json:"circle_id"`
	CreatedAt *time.Time `json:"created_at"`
}

func (UserBookmark) TableName() string {
	return "user_bookmark"
}

type UserUpvote struct {
	UserID    int        `json:"user_id"`
	CircleID  int        `json:"circle_id"`
	CreatedAt *time.Time `json:"created_at"`
}

func (UserUpvote) TableName() string {
	return "user_upvote"
}
