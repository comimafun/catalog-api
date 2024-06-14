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
	circleblock "catalog-be/internal/modules/circle_block"
	"catalog-be/internal/modules/product"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"errors"
	"strings"
)

type CircleService interface {
	OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*circle_dto.CircleResponse, *domain.Error)
	PublishCircleByID(circleID int) (*string, *domain.Error)
	FindCircleBySlug(slug string, userID int) (*circle_dto.CircleResponse, *domain.Error)
	FindCircleByID(circleID int) (*entity.Circle, *domain.Error)
	UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error)

	transformCircleRawToCircleResponse(rows []entity.CircleRaw) []circle_dto.CircleResponse
	GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter, userID int) (*dto.Pagination[[]circle_dto.CircleResponse], *domain.Error)
	GetPaginatedBookmarkedCircle(userID int, filter *circle_dto.FindAllCircleFilter) (*dto.Pagination[[]circle_dto.CircleResponse], *domain.Error)

	SaveBookmarkCircle(circleID int, userID int) *domain.Error
	DeleteBookmarkCircle(circleID int, userID int) *domain.Error
}

type circleService struct {
	circleRepo            CircleRepo
	userService           user.UserService
	utils                 utils.Utils
	refreshTokenService   refreshtoken.RefreshTokenService
	circleBlockService    circleblock.CircleBlockService
	circleWorkTypeService circle_work_type.CircleWorkTypeService
	circleFandomService   circle_fandom.CircleFandomService
	bookmark              bookmark.CircleBookmarkService
	productService        product.ProductService
}

// UpdateCircleByID implements CircleService.
func (c *circleService) UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error) {
	if body.Name != nil && *body.Name == "" {
		return nil, domain.NewError(400, errors.New("CIRCLE_NAME_CANNOT_BE_EMPTY"), nil)
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

	if body.URL != nil && *body.URL != circle.URL {
		circle.URL = *body.URL
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
		circle.Description = body.Description
	}

	if body.Batch != nil && body.Batch != circle.Batch {
		circle.Batch = body.Batch
	}

	if body.Day != nil && body.Day != circle.Day {
		circle.Day = body.Day
	}

	rows, err := c.circleRepo.UpdateCircleAndAllRelation(circleID, circle, body)
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

func NewCircleService(
	circleRepo CircleRepo,
	userService user.UserService,
	utils utils.Utils,
	refreshTokenService refreshtoken.RefreshTokenService,
	circleBlockService circleblock.CircleBlockService,
	circleWorkTypeService circle_work_type.CircleWorkTypeService,
	circleFandomService circle_fandom.CircleFandomService,
	bookmark bookmark.CircleBookmarkService,
	product product.ProductService,
) CircleService {
	return &circleService{
		circleRepo:            circleRepo,
		userService:           userService,
		utils:                 utils,
		refreshTokenService:   refreshTokenService,
		circleBlockService:    circleBlockService,
		circleWorkTypeService: circleWorkTypeService,
		circleFandomService:   circleFandomService,
		bookmark:              bookmark,
		productService:        product,
	}
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

				productExist := false

				for _, product := range response[i].Product {
					if product.ID == row.ProductID {
						productExist = true
						break
					}
				}

				if !productExist && row.ProductID != 0 {
					response[i].Product = append(response[i].Product, entity.Product{
						ID:        row.ProductID,
						Name:      row.ProductName,
						ImageURL:  row.ProductImageURL,
						CircleID:  row.ID,
						CreatedAt: row.ProductCreatedAt,
						UpdatedAt: row.ProductUpdatedAt,
					})
				}

			}

		}

		if !found {
			latestRow := circle_dto.CircleResponse{
				Circle: entity.Circle{
					ID:           row.ID,
					Name:         row.Name,
					Slug:         row.Slug,
					PictureURL:   row.PictureURL,
					FacebookURL:  row.FacebookURL,
					InstagramURL: row.InstagramURL,
					TwitterURL:   row.TwitterURL,
					Description:  row.Description,
					Batch:        row.Batch,
					Verified:     row.Verified,
					Published:    row.Published,
					CreatedAt:    row.CreatedAt,
					UpdatedAt:    row.UpdatedAt,
					DeletedAt:    row.DeletedAt,
					Day:          row.Day,
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

			if row.ProductID != 0 {
				latestRow.Product = append(latestRow.Product, entity.Product{
					ID:        row.ProductID,
					Name:      row.ProductName,
					ImageURL:  row.ProductImageURL,
					CircleID:  row.ID,
					CreatedAt: row.ProductCreatedAt,
					UpdatedAt: row.ProductUpdatedAt,
				})
			} else {
				latestRow.Product = []entity.Product{}
			}

			if row.BlockID != 0 {
				latestRow.Block = &circle_dto.BlockResponse{
					ID:        row.BlockID,
					Name:      row.BlockName,
					CreatedAt: row.BlockCreatedAt,
					UpdatedAt: row.BlockUpdatedAt,
				}
			}

			response = append(response, latestRow)
		}
	}

	return response
}

// GetPaginatedBookmarkedCircle implements CircleService.
func (c *circleService) GetPaginatedBookmarkedCircle(userID int, filter *circle_dto.FindAllCircleFilter) (*dto.Pagination[[]circle_dto.CircleResponse], *domain.Error) {
	rows, err := c.circleRepo.FindBookmarkedCircleByUserID(userID, filter)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleResponse(rows)

	count, err := c.circleRepo.FindAllBookmarkedCount(userID, filter)
	if err != nil {
		return nil, err
	}
	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	return &dto.Pagination[[]circle_dto.CircleResponse]{
		Data:     response,
		Metadata: *metadata,
	}, nil
}

// GetPaginatedCircle implements CircleService.
func (c *circleService) GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter, userID int) (*dto.Pagination[[]circle_dto.CircleResponse], *domain.Error) {
	rows, err := c.circleRepo.FindAllCircles(filter, userID)
	if err != nil {
		return nil, err
	}

	response := c.transformCircleRawToCircleResponse(rows)

	count, countErr := c.circleRepo.FindAllCount(filter)
	if countErr != nil {
		return nil, countErr
	}

	metadata := factory.GetPaginationMetadata(count, filter.Page, filter.Limit)

	return &dto.Pagination[[]circle_dto.CircleResponse]{
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
	slug = slug + "-" + c.utils.GenerateRandomCode(2)
	circle, err := c.circleRepo.CreateOne(entity.Circle{
		Name:         body.Name,
		Slug:         slug,
		PictureURL:   &body.PictureURL,
		FacebookURL:  &body.FacebookURL,
		InstagramURL: &body.InstagramURL,
		TwitterURL:   &body.TwitterURL,
	})

	if err != nil {
		return nil, err
	}

	_, updateErr := c.userService.UpdateOneByID(entity.User{
		ID:       userID,
		CircleID: &circle.ID,
	})
	if updateErr != nil {
		return nil, updateErr
	}

	return &circle_dto.CircleResponse{
		Circle:   *circle,
		Fandom:   []entity.Fandom{},
		WorkType: []entity.WorkType{},
		Product:  []entity.Product{},
	}, nil

}
