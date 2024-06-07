package circle

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	circle_dto "catalog-be/internal/modules/circle/dto"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"errors"
)

type CircleService interface {
	OnboardNewCircle(body *circle_dto.OnboardNewCircleRequestBody, userID int) (*entity.Circle, *domain.Error)
	PublishCircleByID(circleID int) (*string, *domain.Error)
	FindCircleBySlug(slug string) (*entity.Circle, *domain.Error)
}

type circleService struct {
	circleRepo          CircleRepo
	userService         user.UserService
	utils               utils.Utils
	refreshTokenService refreshtoken.RefreshTokenService
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
) CircleService {
	return &circleService{
		circleRepo:          circleRepo,
		userService:         userService,
		utils:               utils,
		refreshTokenService: refreshTokenService,
	}
}
