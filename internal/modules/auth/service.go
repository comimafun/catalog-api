package auth

import (
	internal_config "catalog-be/internal/config"
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	auth_dto "catalog-be/internal/modules/auth/dto"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService interface {
	AuthWithGoogle(code string) (*auth_dto.NewTokenResponse, *domain.Error)
	GetAuthURL() string
	login(user *entity.User) (*auth_dto.NewTokenResponse, *domain.Error)
	registerWithGoogle(user *auth_dto.GoogleUserData) (*entity.User, *domain.Error)
	generateNewJWTAndRefreshToken(user *entity.User) (*auth_dto.NewToken, *domain.Error)
}

type authService struct {
	userService         user.UserService
	config              internal_config.Config
	refreshTokenService refreshtoken.RefreshTokenService
	utils               utils.Utils
}

// GetAuthURL implements AuthService.
func (a *authService) GetAuthURL() string {
	return a.config.AuthCodeURL()
}

// registerWithGoogle implements AuthService.
func (a *authService) registerWithGoogle(user *auth_dto.GoogleUserData) (*entity.User, *domain.Error) {
	randString := a.utils.GenerateRandomCode(10)
	hash, hashingErr := a.utils.HashPassword(randString)
	if hashingErr != nil {
		return nil, hashingErr
	}
	newUser, err := a.userService.CreateOne(entity.User{
		Name:              user.Name,
		Email:             user.Email,
		ProfilePictureURL: user.Picture,
		Hash:              *hash,
	})
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

// generateNewJWTAndRefreshToken implements AuthService.
func (a *authService) generateNewJWTAndRefreshToken(user *entity.User) (*auth_dto.NewToken, *domain.Error) {
	secret := os.Getenv("JWT_SECRET")
	expiredAt := time.Now().Add(time.Minute * 60)
	claims := auth_dto.ATClaims{
		UserID:   user.ID,
		Email:    user.Email,
		CircleID: user.CircleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, signErr := token.SignedString([]byte(secret))

	if signErr != nil {
		return nil, domain.NewError(500, signErr, nil)
	}

	refreshToken := utils.NewUtils().GenerateRandomCode(20)
	refreshTokenExpiredAt := time.Now().Add(time.Hour * 24 * 30)

	return &auth_dto.NewToken{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  expiredAt,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
	}, nil
}

// login implements AuthService.
func (a *authService) login(user *entity.User) (*auth_dto.NewTokenResponse, *domain.Error) {
	newToken, newTokenErr := a.generateNewJWTAndRefreshToken(user)
	if newTokenErr != nil {
		return nil, newTokenErr
	}

	_, insertErr := a.refreshTokenService.CreateOne(entity.RefreshToken{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		UserID:       user.ID,
		ExpiredAt:    &newToken.RefreshTokenExpiredAt,
	})
	if insertErr != nil {
		return nil, insertErr
	}

	return &auth_dto.NewTokenResponse{
		AccessToken:           newToken.AccessToken,
		RefreshToken:          newToken.RefreshToken,
		AccessTokenExpiredAt:  newToken.AccessTokenExpiredAt.Format(time.RFC3339),
		RefreshTokenExpiredAt: newToken.RefreshTokenExpiredAt.Format(time.RFC3339),
	}, nil
}

// AuthWithGoogle implements AuthService.
func (a *authService) AuthWithGoogle(code string) (*auth_dto.NewTokenResponse, *domain.Error) {
	user, err := a.config.ParseCodeToUserData(code)
	if err != nil {
		return nil, err
	}

	existingUser, existingUserErr := a.userService.FindOneByEmail(user.Email)

	if existingUserErr != nil && !errors.Is(existingUserErr.Err, gorm.ErrRecordNotFound) {
		return nil, existingUserErr
	}

	if existingUser == nil {
		login, loginErr := a.login(existingUser)
		if loginErr != nil {
			return nil, loginErr
		}
		return login, nil
	}

	newUser, newUserErr := a.registerWithGoogle(user)
	if newUserErr != nil {
		return nil, newUserErr
	}

	login, loginErr := a.login(newUser)

	if loginErr != nil {
		return nil, loginErr
	}

	return login, nil

}

func NewAuthService(
	userService user.UserService,
	config internal_config.Config,
	refreshToken refreshtoken.RefreshTokenService,
	utils utils.Utils,
) AuthService {
	return &authService{
		userService,
		config,
		refreshToken,
		utils,
	}
}
