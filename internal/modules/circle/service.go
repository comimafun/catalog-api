package circle

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
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
	FindCircleBySlug(slug string) (*entity.Circle, *domain.Error)
	UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*circle_dto.CircleResponse, *domain.Error)
	GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter) (*[]circle_dto.CircleResponse, *domain.Error)
}

type circleService struct {
	circleRepo          CircleRepo
	userService         user.UserService
	utils               utils.Utils
	refreshTokenService refreshtoken.RefreshTokenService
	circleBlockService  circleblock.CircleBlockService
}

// GetPaginatedCircle implements CircleService.
func (c *circleService) GetPaginatedCircle(filter *circle_dto.FindAllCircleFilter) (*[]circle_dto.CircleResponse, *domain.Error) {
	circles, err := c.circleRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	var circleMap = make(map[int]circle_dto.CircleResponse)

	for _, circle := range circles {
		if _, ok := circleMap[circle.ID]; !ok {
			circleMap[circle.ID] = circle_dto.CircleResponse{
				Circle: circle.Circle,
				Fandom: []entity.Fandom{},
			}
		}

		if circle.FandomID != 0 {
			fandom := entity.Fandom{
				ID:        circle.FandomID,
				Name:      circle.FandomName,
				Visible:   circle.FandomVisible,
				CreatedAt: circle.FandomCreatedAt,
				UpdatedAt: circle.FandomUpdatedAt,
				DeletedAt: circle.FandomDeletedAt,
			}

			circleResponse := circleMap[circle.ID]
			circleResponse.Fandom = append(circleResponse.Fandom, fandom)
			circleMap[circle.ID] = circleResponse
		}
	}

	var response []circle_dto.CircleResponse

	for _, circle := range circleMap {
		response = append(response, circle)

	}

	return &response, nil
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

	if len(body.Fandom) > 0 {
		for _, fandom := range body.Fandom {
			err := c.circleRepo.UpsertCircleFandomRelation(circleID, fandom.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	fandoms, fandomErr := c.circleRepo.FindAllCircleRelationFandom(circleID)
	if fandomErr != nil {
		return nil, fandomErr
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
		Circle: *updated,
		Fandom: fandoms,
	}, nil
}

// FindCircleBySlug implements CircleService.
func (c *circleService) FindCircleBySlug(slug string) (*entity.Circle, *domain.Error) {
	if slug == "" {
		return nil, domain.NewError(400, errors.New("SLUG_IS_EMPTY"), nil)
	}

	circle, err := c.circleRepo.FindOneBySlug(slug)
	if err != nil {
		return nil, err
	}

	return circle, nil
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

func NewCircleService(
	circleRepo CircleRepo,
	userService user.UserService,
	utils utils.Utils,
	refreshTokenService refreshtoken.RefreshTokenService,
	circleBlockService circleblock.CircleBlockService,
) CircleService {
	return &circleService{
		circleRepo:          circleRepo,
		userService:         userService,
		utils:               utils,
		refreshTokenService: refreshTokenService,
		circleBlockService:  circleBlockService,
	}
}
