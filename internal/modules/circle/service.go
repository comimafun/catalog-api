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
	UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*entity.Circle, *domain.Error)
}

type circleService struct {
	circleRepo          CircleRepo
	userService         user.UserService
	utils               utils.Utils
	refreshTokenService refreshtoken.RefreshTokenService
	circleBlockService  circleblock.CircleBlockService
}

// UpdateCircleByID implements CircleService.
func (c *circleService) UpdateCircleByID(circleID int, body *circle_dto.UpdateCircleRequestBody) (*entity.Circle, *domain.Error) {
	if body.CircleBlock != nil {
		trimmedBlock := strings.TrimSpace(*body.CircleBlock)
		if trimmedBlock == "" {
			return nil, domain.NewError(400, errors.New("CIRCLE_BLOCK_IS_EMPTY"), nil)
		}

		body.CircleBlock = &trimmedBlock

		block, err := c.circleBlockService.GetOneByBlock(*body.CircleBlock)
		if err != nil && !errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if block != nil {
			return nil, domain.NewError(400, errors.New("CIRCLE_BLOCK_ALREADY_EXIST"), nil)
		}

		_, newCircleBlockErr := c.circleBlockService.CreateOne(*body.CircleBlock, circleID)
		if newCircleBlockErr != nil {
			return nil, newCircleBlockErr
		}
	}

	updated, err := c.circleRepo.UpdateOneByID(circleID, entity.Circle{
		PictureURL:   body.PictureURL,
		FacebookURL:  body.FacebookURL,
		InstagramURL: body.InstagramURL,
		TwitterURL:   body.TwitterURL,
		Day:          body.Day,
		Description:  body.Description,
		Batch:        body.Batch,
	})

	if err != nil {
		return nil, err
	}

	return updated, nil
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
		PictureURL:   body.PictureURL,
		FacebookURL:  body.FacebookURL,
		InstagramURL: body.InstagramURL,
		TwitterURL:   body.TwitterURL,
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
