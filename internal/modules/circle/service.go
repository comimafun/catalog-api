package circle

import (
	"catalog-be/internal/database/factory"
	"catalog-be/internal/domain"
	"catalog-be/internal/dto"
	"catalog-be/internal/entity"
	"catalog-be/internal/modules/circle/bookmark"
	"catalog-be/internal/modules/circle/circle_fandom"
	"catalog-be/internal/modules/circle/circle_work_type"
	circle_dto "catalog-be/internal/modules/circle/dto"
	"catalog-be/internal/modules/product"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"catalog-be/internal/validation"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CircleService interface {
	OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*circle_dto.CircleResponse, *domain.Error)
	PublishCircleByID(circleID int) (*string, *domain.Error)
	FindCircleBySlug(slug string, userID int) (*circle_dto.CircleResponse, *domain.Error)
	FindCircleByID(circleID int) (*entity.Circle, *domain.Error)
	UpdateCircleByID(userID int, circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error)

	transformCircleRawToCircleResponse(rows []entity.CircleRaw) []circle_dto.CircleResponse
	transformCircleRawToCircleOneForPaginationResponse(rows []entity.CircleRaw) []circle_dto.CircleOneForPaginationResponse

	GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter, userID int) (*dto.Pagination[[]circle_dto.CircleOneForPaginationResponse], *domain.Error)
	GetPaginatedBookmarkedCircle(userID int, filter *circle_dto.FindAllCircleFilter) (*dto.Pagination[[]circle_dto.CircleOneForPaginationResponse], *domain.Error)

	SaveBookmarkCircle(circleID int, userID int) *domain.Error
	DeleteBookmarkCircle(circleID int, userID int) *domain.Error
}

type circleService struct {
	circleRepo            CircleRepo
	userService           user.UserService
	utils                 utils.Utils
	refreshTokenService   refreshtoken.RefreshTokenService
	circleWorkTypeService circle_work_type.CircleWorkTypeService
	circleFandomService   circle_fandom.CircleFandomService
	bookmark              bookmark.CircleBookmarkService
	productService        product.ProductService
	sanitizer             *validation.Sanitizer
}

func NewCircleService(
	circleRepo CircleRepo,
	userService user.UserService,
	utils utils.Utils,
	refreshTokenService refreshtoken.RefreshTokenService,
	circleWorkTypeService circle_work_type.CircleWorkTypeService,
	circleFandomService circle_fandom.CircleFandomService,
	bookmark bookmark.CircleBookmarkService,
	product product.ProductService,
	sanitizer *validation.Sanitizer,
) CircleService {
	return &circleService{
		circleRepo:            circleRepo,
		userService:           userService,
		utils:                 utils,
		refreshTokenService:   refreshTokenService,
		circleWorkTypeService: circleWorkTypeService,
		circleFandomService:   circleFandomService,
		bookmark:              bookmark,
		productService:        product,
		sanitizer:             sanitizer,
	}
}

// transformCircleRawToCircleOneForPaginationResponse implements CircleService.
func (c *circleService) transformCircleRawToCircleOneForPaginationResponse(rows []entity.CircleRaw) []circle_dto.CircleOneForPaginationResponse {
	if len(rows) == 0 {
		return []circle_dto.CircleOneForPaginationResponse{}
	}

	var response []circle_dto.CircleOneForPaginationResponse

	for _, row := range rows {
		// check if row is inside response
		var found bool
		for i, res := range response {
			if res.ID == row.ID {
				found = true

				fandomExist := false

				for _, fandom := range response[i].Fandom {
					if fandom.ID == row.FandomID {
						fandomExist = true
						break
					}
				}

				if !fandomExist && row.FandomID != 0 {
					response[i].Fandom = append(response[i].Fandom, entity.Fandom{
						ID:        row.FandomID,
						Name:      row.FandomName,
						CreatedAt: row.FandomCreatedAt,
						UpdatedAt: row.FandomUpdatedAt,
					})
				}

				workTypeExist := false

				for _, workType := range response[i].WorkType {
					if workType.ID == row.WorkTypeID {
						workTypeExist = true
						break
					}
				}

				if !workTypeExist && row.WorkTypeID != 0 {
					response[i].WorkType = append(response[i].WorkType, entity.WorkType{
						ID:        row.WorkTypeID,
						Name:      row.WorkTypeName,
						CreatedAt: row.WorkTypeCreatedAt,
						UpdatedAt: row.WorkTypeUpdatedAt,
					})
				}

			}

		}

		if !found {
			latestRow := circle_dto.CircleOneForPaginationResponse{
				Circle: entity.Circle{
					ID:              row.ID,
					Name:            row.Name,
					Slug:            row.Slug,
					PictureURL:      row.PictureURL,
					FacebookURL:     row.FacebookURL,
					InstagramURL:    row.InstagramURL,
					TwitterURL:      row.TwitterURL,
					Description:     row.Description,
					Batch:           row.Batch,
					Verified:        row.Verified,
					Published:       row.Published,
					CreatedAt:       row.CreatedAt,
					UpdatedAt:       row.UpdatedAt,
					DeletedAt:       row.DeletedAt,
					Day:             row.Day,
					URL:             row.URL,
					EventID:         row.EventID,
					CoverPictureURL: row.CoverPictureURL,
				},
				Fandom:     []entity.Fandom{},
				WorkType:   []entity.WorkType{},
				Bookmarked: row.Bookmarked,
			}

			if row.FandomID != 0 {
				latestRow.Fandom = append(latestRow.Fandom, entity.Fandom{
					ID:        row.FandomID,
					Name:      row.FandomName,
					CreatedAt: row.FandomCreatedAt,
					UpdatedAt: row.FandomUpdatedAt,
				})
			} else {
				latestRow.Fandom = []entity.Fandom{}
			}

			if row.WorkTypeID != 0 {
				latestRow.WorkType = append(latestRow.WorkType, entity.WorkType{
					ID:        row.WorkTypeID,
					Name:      row.WorkTypeName,
					CreatedAt: row.WorkTypeCreatedAt,
					UpdatedAt: row.WorkTypeUpdatedAt,
				})
			} else {
				latestRow.WorkType = []entity.WorkType{}
			}

			if row.BlockEventID != 0 {
				latestRow.BlockEvent = &circle_dto.BlockResponse{
					ID:   row.BlockEventID,
					Name: row.BlockEventName,
				}
			}

			response = append(response, latestRow)
		}
	}

	return response
}

// UpdateCircleByID implements CircleService.
func (c *circleService) UpdateCircleByID(userID int, circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error) {
	if body.Name != nil && *body.Name == "" {
		return nil, domain.NewError(400, errors.New("CIRCLE_NAME_CANNOT_BE_EMPTY"), nil)
	}

	if body.EventID == nil {
		if body.Day != nil || body.CircleBlock != nil {
			return nil, domain.NewError(400, errors.New("EVENT_ID_CANNOT_BE_EMPTY"), nil)
		}

	}

	circle, err := c.FindCircleByID(circleID)
	if err != nil {
		return nil, err
	}

	if body.Name != nil && *body.Name != circle.Name {
		slug, slugErr := c.utils.Slugify(*body.Name)
		if slugErr != nil {
			return nil, domain.NewError(500, slugErr, nil)
		}
		circle.Name = *body.Name
		circle.Slug = strings.ToLower(slug + "-" + c.utils.GenerateRandomCode(2))
	}

	if body.URL != nil && body.URL != circle.URL {
		circle.URL = body.URL
	}

	if body.PictureURL != nil && *body.PictureURL != *circle.PictureURL {
		circle.PictureURL = body.PictureURL
	}

	if body.FacebookURL != nil && *body.FacebookURL != *circle.FacebookURL {
		circle.FacebookURL = body.FacebookURL
	}

	if body.InstagramURL != nil && *body.InstagramURL != *circle.InstagramURL {
		circle.InstagramURL = body.InstagramURL
	}

	if body.TwitterURL != nil && *body.TwitterURL != *circle.TwitterURL {
		circle.TwitterURL = body.TwitterURL
	}

	if body.Description != nil && body.Description != circle.Description {
		sanitized := c.sanitizer.Sanitize(*body.Description)
		circle.Description = &sanitized
	}

	if body.Batch != nil && body.Batch != circle.Batch {
		circle.Batch = body.Batch
	}

	if body.Day != nil && body.Day != circle.Day {
		if *body.Day == "" {
			circle.Day = nil
		} else {
			circle.Day = body.Day
		}
	}

	if body.EventID != nil && body.EventID != circle.EventID {
		if *body.EventID == 0 {
			circle.EventID = nil
		} else {
			circle.EventID = body.EventID
		}
	}

	if body.CircleBlock != nil {
		trimmedBlock := strings.TrimSpace(*body.CircleBlock)
		if trimmedBlock != "" && body.EventID == nil {
			return nil, domain.NewError(400, errors.New("EVENT_ID_CANNOT_BE_EMPTY"), nil)
		}

		if trimmedBlock != "" {
			body.CircleBlock = &trimmedBlock
		}
	}

	if body.CoverPictureURL != nil && *body.CoverPictureURL != *circle.CoverPictureURL {
		if *body.CoverPictureURL == "" {
			circle.CoverPictureURL = nil
		} else {
			circle.CoverPictureURL = body.CoverPictureURL
		}
	}

	rows, err := c.circleRepo.UpdateCircleAndAllRelation(userID, circle, body)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleResponse(rows)

	if len(response) == 0 {
		return nil, domain.NewError(404, errors.New("CIRCLE_NOT_FOUND"), nil)
	}

	return &response[0], nil
}

// FindCircleByID implements CircleService.
func (c *circleService) FindCircleByID(circleID int) (*entity.Circle, *domain.Error) {
	circle, err := c.circleRepo.FindOneByID(circleID)
	if err != nil {
		return nil, err
	}

	return circle, nil
}

// DeleteBookmarkCircle implements CircleService.
func (c *circleService) DeleteBookmarkCircle(circleID int, userID int) *domain.Error {
	return c.bookmark.DeleteBookmark(circleID, userID)
}

// SaveBookmarkCircle implements CircleService.
func (c *circleService) SaveBookmarkCircle(circleID int, userID int) *domain.Error {
	return c.bookmark.CreateBookmark(circleID, userID)
}

// transformCircleRawToCircleResponse implements CircleService.
func (c *circleService) transformCircleRawToCircleResponse(rows []entity.CircleRaw) []circle_dto.CircleResponse {
	if len(rows) == 0 {
		return []circle_dto.CircleResponse{}
	}

	var response []circle_dto.CircleResponse

	for _, row := range rows {
		// check if row is inside response
		var found bool
		for i, res := range response {
			if res.ID == row.ID {
				found = true

				fandomExist := false

				for _, fandom := range response[i].Fandom {
					if fandom.ID == row.FandomID {
						fandomExist = true
						break
					}
				}

				if !fandomExist && row.FandomID != 0 {
					response[i].Fandom = append(response[i].Fandom, entity.Fandom{
						ID:        row.FandomID,
						Name:      row.FandomName,
						CreatedAt: row.FandomCreatedAt,
						UpdatedAt: row.FandomUpdatedAt,
					})
				}

				workTypeExist := false

				for _, workType := range response[i].WorkType {
					if workType.ID == row.WorkTypeID {
						workTypeExist = true
						break
					}
				}

				if !workTypeExist && row.WorkTypeID != 0 {
					response[i].WorkType = append(response[i].WorkType, entity.WorkType{
						ID:        row.WorkTypeID,
						Name:      row.WorkTypeName,
						CreatedAt: row.WorkTypeCreatedAt,
						UpdatedAt: row.WorkTypeUpdatedAt,
					})
				}

			}

		}

		if !found {
			latestRow := circle_dto.CircleResponse{
				Circle: entity.Circle{
					ID:              row.ID,
					Name:            row.Name,
					Slug:            row.Slug,
					PictureURL:      row.PictureURL,
					FacebookURL:     row.FacebookURL,
					InstagramURL:    row.InstagramURL,
					TwitterURL:      row.TwitterURL,
					Description:     row.Description,
					Batch:           row.Batch,
					Verified:        row.Verified,
					Published:       row.Published,
					CreatedAt:       row.CreatedAt,
					UpdatedAt:       row.UpdatedAt,
					DeletedAt:       row.DeletedAt,
					Day:             row.Day,
					URL:             row.URL,
					EventID:         row.EventID,
					CoverPictureURL: row.CoverPictureURL,
				},
				Fandom:     []entity.Fandom{},
				WorkType:   []entity.WorkType{},
				Bookmarked: row.Bookmarked,
			}

			if row.FandomID != 0 {
				latestRow.Fandom = append(latestRow.Fandom, entity.Fandom{
					ID:        row.FandomID,
					Name:      row.FandomName,
					CreatedAt: row.FandomCreatedAt,
					UpdatedAt: row.FandomUpdatedAt,
				})
			} else {
				latestRow.Fandom = []entity.Fandom{}
			}

			if row.WorkTypeID != 0 {
				latestRow.WorkType = append(latestRow.WorkType, entity.WorkType{
					ID:        row.WorkTypeID,
					Name:      row.WorkTypeName,
					CreatedAt: row.WorkTypeCreatedAt,
					UpdatedAt: row.WorkTypeUpdatedAt,
				})
			} else {
				latestRow.WorkType = []entity.WorkType{}
			}

			if row.EventID != nil {
				latestRow.Event = &entity.Event{
					ID:          *row.EventID,
					Name:        row.EventName,
					Slug:        row.EventSlug,
					Description: row.EventDescription,
					StartedAt:   *row.EventStartedAt,
					EndedAt:     *row.EventEndedAt,
				}
			}

			if row.BlockEventID != 0 {
				latestRow.BlockEvent = &circle_dto.BlockResponse{
					ID:   row.BlockEventID,
					Name: row.BlockEventName,
				}

			}

			response = append(response, latestRow)
		}
	}

	return response
}

// GetPaginatedBookmarkedCircle implements CircleService.
func (c *circleService) GetPaginatedBookmarkedCircle(userID int, filter *circle_dto.FindAllCircleFilter) (*dto.Pagination[[]circle_dto.CircleOneForPaginationResponse], *domain.Error) {
	rows, err := c.circleRepo.FindBookmarkedCircleByUserID(userID, filter)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleOneForPaginationResponse(rows)

	count, err := c.circleRepo.FindAllBookmarkedCount(userID, filter)
	if err != nil {
		return nil, err
	}
	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	return &dto.Pagination[[]circle_dto.CircleOneForPaginationResponse]{
		Data:     response,
		Metadata: *metadata,
	}, nil
}

// GetPaginatedCircle implements CircleService.
func (c *circleService) GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter, userID int) (*dto.Pagination[[]circle_dto.CircleOneForPaginationResponse], *domain.Error) {
	rows, err := c.circleRepo.FindAllCircles(filter, userID)
	if err != nil {
		return nil, err
	}
	count, err := c.circleRepo.FindAllCount(filter)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleOneForPaginationResponse(rows)
	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	return &dto.Pagination[[]circle_dto.CircleOneForPaginationResponse]{
		Data:     response,
		Metadata: *metadata,
	}, nil

}

// FindCircleBySlug implements CircleService.
func (c *circleService) FindCircleBySlug(slug string, userID int) (*circle_dto.CircleResponse, *domain.Error) {
	if slug == "" {
		return nil, domain.NewError(400, errors.New("SLUG_IS_EMPTY"), nil)
	}

	rows, err := c.circleRepo.FindOneBySlugAndRelatedTables(slug, userID)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleResponse(rows)

	if len(response) == 0 {
		return nil, domain.NewError(404, errors.New("CIRCLE_NOT_FOUND"), nil)
	}

	return &response[0], nil
}

// PublishCircleByID implements CircleService.
func (c *circleService) PublishCircleByID(circleID int) (*string, *domain.Error) {
	if circleID == 0 {
		return nil, domain.NewError(400, errors.New("CIRCLE_NOT_FOUND"), nil)
	}

	circle, err := c.circleRepo.FindOneByID(circleID)
	if err != nil {
		return nil, err
	}

	circle.Published = !circle.Published

	_, err = c.circleRepo.UpserstOneCircle(circle)

	if err != nil {
		return nil, err
	}
	published := "PUBLISHED"
	if circle.Published {
		published = "UNPUBLISHED"
	}

	return &published, nil
}

// OnboardNewCircle implements CircleService.
func (c *circleService) OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*circle_dto.CircleResponse, *domain.Error) {
	slug, slugErr := c.utils.Slugify(body.Name)
	if slugErr != nil {
		return nil, domain.NewError(500, slugErr, nil)
	}

	user, userErr := c.userService.FindOneByID(userID)
	if userErr != nil {
		if errors.Is(userErr.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, errors.New("USER_NOT_FOUND"), nil)
		}
		return nil, userErr
	}

	slug = slug + "-" + c.utils.GenerateRandomCode(2)
	circle, err := c.circleRepo.OnboardNewCircle(&entity.Circle{
		Name:         body.Name,
		Slug:         strings.ToLower(slug),
		PictureURL:   &body.PictureURL,
		FacebookURL:  &body.FacebookURL,
		InstagramURL: &body.InstagramURL,
		TwitterURL:   &body.TwitterURL,
		URL:          &body.URL,
	}, user)

	if err != nil {
		return nil, err
	}

	return &circle_dto.CircleResponse{
		Circle:   *circle,
		Fandom:   []entity.Fandom{},
		WorkType: []entity.WorkType{},
	}, nil

}
