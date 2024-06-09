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
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CircleService interface {
	OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*entity.Circle, *domain.Error)
	PublishCircleByID(circleID int) (*string, *domain.Error)
	FindCircleBySlug(slug string, userID int) (*circle_dto.CircleResponse, *domain.Error)
	UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error)

	transformCircleRawToCircleResponse(rows []entity.CircleRaw) ([]circle_dto.CircleResponse, *domain.Error)
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
func (c *circleService) transformCircleRawToCircleResponse(rows []entity.CircleRaw) ([]circle_dto.CircleResponse, *domain.Error) {
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
				Fandom:       []entity.Fandom{},
				WorkType:     []entity.WorkType{},
				Bookmarked:   row.Bookmarked,
				BookmarkedAt: row.BookmarkedAt,
			}

			if row.FandomID != 0 {
				latestRow.Fandom = append(latestRow.Fandom, entity.Fandom{
					ID:        row.FandomID,
					Name:      row.FandomName,
					CreatedAt: row.FandomCreatedAt,
					UpdatedAt: row.FandomUpdatedAt,
				})
			}

			if row.WorkTypeID != 0 {
				latestRow.WorkType = append(latestRow.WorkType, entity.WorkType{
					ID:        row.WorkTypeID,
					Name:      row.WorkTypeName,
					CreatedAt: row.WorkTypeCreatedAt,
					UpdatedAt: row.WorkTypeUpdatedAt,
				})
			}

			if row.BlockID != 0 {
				latestRow.Block = &circle_dto.BlockResponse{
					ID:        row.BlockID,
					Block:     row.BlockPrefix + "-" + row.BlockPostfix,
					CreatedAt: row.BlockCreatedAt,
					UpdatedAt: row.BlockUpdatedAt,
				}
			}

			response = append(response, latestRow)
		}
	}

	return response, nil
}

// GetPaginatedBookmarkedCircle implements CircleService.
func (c *circleService) GetPaginatedBookmarkedCircle(userID int, filter *circle_dto.FindAllCircleFilter) (*dto.Pagination[[]circle_dto.CircleResponse], *domain.Error) {
	rows, err := c.circleRepo.FindBookmarkedCircleByUserID(userID, filter)
	if err != nil {
		return nil, err
	}

	response, err := c.transformCircleRawToCircleResponse(rows)

	if err != nil {
		return nil, err
	}

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

	response, err := c.transformCircleRawToCircleResponse(rows)

	if err != nil {
		return nil, err
	}

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

// UpdateCircleByID implements CircleService.
func (c *circleService) UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error) {
	trimmedBlock := strings.TrimSpace(body.CircleBlock)
	if trimmedBlock != "" {
		body.CircleBlock = trimmedBlock

		block, err := c.circleBlockService.GetOneByBlock(body.CircleBlock)
		if err != nil && !errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if block != nil {
			return nil, domain.NewError(400, errors.New("CIRCLE_BLOCK_ALREADY_EXIST"), nil)
		}

		_, newCircleBlockErr := c.circleBlockService.CreateOne(body.CircleBlock, circleID)
		if newCircleBlockErr != nil {
			return nil, newCircleBlockErr
		}
	}

	if len(body.FandomIDs) > 0 {
		if len(body.FandomIDs) > 5 {
			return nil, domain.NewError(400, errors.New("FANDOM_LIMIT_EXCEEDED"), nil)
		}

		insertErr := c.circleFandomService.BulkDeleteAndInsertCircleFandomType(circleID, body.FandomIDs)
		if insertErr != nil {
			return nil, insertErr
		}
	}

	if len(body.WorkTypeIDs) > 0 {
		if len(body.WorkTypeIDs) > 5 {
			return nil, domain.NewError(400, errors.New("WORK_TYPE_LIMIT_EXCEEDED"), nil)
		}

		insertErr := c.circleWorkTypeService.BulkDeleteAndInsertCircleWorkType(circleID, body.WorkTypeIDs)
		if insertErr != nil {
			return nil, insertErr
		}
	}

	fandoms, fandomErr := c.circleFandomService.FindAllCircleFandomTypeRelated(circleID)
	if fandomErr != nil {
		return nil, fandomErr
	}

	workType, workTypeErr := c.circleWorkTypeService.FindAllCircleWorkTypeRelated(circleID)
	if workTypeErr != nil {
		return nil, workTypeErr
	}

	updated, err := c.circleRepo.UpdateOneByID(circleID, entity.Circle{
		PictureURL:   &body.PictureURL,
		FacebookURL:  &body.FacebookURL,
		InstagramURL: &body.InstagramURL,
		TwitterURL:   &body.TwitterURL,
		Day:          body.Day,
		Description:  &body.Description,
		Batch:        body.Batch,
	})

	if err != nil {
		return nil, err
	}

	return &circle_dto.CircleResponse{
		Circle:   *updated,
		Fandom:   fandoms,
		WorkType: workType,
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

	response, err := c.transformCircleRawToCircleResponse(rows)
	if err != nil {
		return nil, err
	}

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

	_, err = c.circleRepo.UpdateOneByID(circleID, entity.Circle{
		Published: !circle.Published,
	})

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
func (c *circleService) OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*entity.Circle, *domain.Error) {
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

	return circle, nil

}
