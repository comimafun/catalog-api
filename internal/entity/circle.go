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
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Slug            string         `json:"slug"`
	URL             *string        `json:"url"`
	PictureURL      *string        `json:"picture_url"`
	CoverPictureURL *string        `json:"cover_picture_url"`
	FacebookURL     *string        `json:"facebook_url"`
	InstagramURL    *string        `json:"instagram_url"`
	TwitterURL      *string        `json:"twitter_url"`
	Description     *string        `json:"description"`
	Rating          *string        `json:"rating"` // enum GA, PG, M
	Verified        bool           `json:"verified"`
	Published       bool           `json:"published"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`

	Day *Day `json:"day"`

	EventID            *int `json:"event_id"`
	UsedReferralCodeID *int `json:"-"`
}

type CircleJoinedTables struct {
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

type CircleRawNew struct {
	Circle

	FandomID        int            `json:"fandom.id" gorm:"column:fandom.id"`
	FandomName      string         `json:"fandom.name" gorm:"column:fandom.name"`
	FandomVisible   bool           `json:"fandom.visible" gorm:"column:fandom.visible"`
	FandomCreatedAt *time.Time     `json:"fandom.created_at"  gorm:"column:fandom.created_at"`
	FandomUpdatedAt *time.Time     `json:"fandom.updated_at" gorm:"column:fandom.updated_at"`
	FandomDeletedAt gorm.DeletedAt `json:"-"`

	WorkTypeID        int            `json:"work_type.id" gorm:"column:work_type.id"`
	WorkTypeName      string         `json:"work_type.name" gorm:"column:work_type.name"`
	WorkTypeCreatedAt *time.Time     `json:"work_type.created_at" gorm:"column:work_type.created_at"`
	WorkTypeUpdatedAt *time.Time     `json:"work_type.updated_at" gorm:"column:work_type.updated_at"`
	WorkTypeDeletedAt gorm.DeletedAt `json:"-"`

	Bookmarked   bool       `json:"bookmarked"`
	BookmarkedAt *time.Time `json:"bookmarked_at"`

	ProductID        int        `json:"product.id" gorm:"column:product.id"`
	ProductName      string     `json:"product.name" gorm:"column:product.name"`
	ProductImageURL  string     `json:"product.image_url" gorm:"column:product.image_url"`
	ProductCreatedAt *time.Time `json:"product.created_at" gorm:"column:product.created_at"`
	ProductUpdatedAt *time.Time `json:"product.updated_at" gorm:"column:product.updated_at"`

	EventName        string     `json:"event.name" gorm:"column:event.name"`
	EventSlug        string     `json:"event.slug" gorm:"column:event.slug"`
	EventDescription string     `json:"event.description" gorm:"column:event.description"`
	EventStartedAt   *time.Time `json:"event.started_at" gorm:"column:event.started_at"`
	EventEndedAt     *time.Time `json:"event.ended_at" gorm:"column:event.ended_at"`

	BlockEventID      int    `json:"block_event.id" gorm:"column:block_event.id"`
	BlockEventPrefix  string `json:"block_event.prefix" gorm:"column:block_event.prefix"`
	BlockEventPostfix string `json:"block_event.postfix" gorm:"column:block_event.postfix"`
	BlockEventName    string `json:"block_event.name" gorm:"column:block_event.name"`
}
