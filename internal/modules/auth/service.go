package auth

import (
	internal_config "catalog-be/internal/config"
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	auth_dto "catalog-be/internal/modules/auth/dto"
	"catalog-be/internal/modules/circle"
	refreshtoken "catalog-be/internal/modules/refresh_token"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/utils"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	userService         *user.UserService
	config              internal_config.Config
	refreshTokenService *refreshtoken.RefreshTokenService
	utils               utils.Utils
	circleService       *circle.CircleService
}

// logoutByAccessToken implements AuthService.
func (a *AuthService) logoutByAccessToken(userID int) *domain.Error {
	err := a.refreshTokenService.DeleteAllRefreshTokenRecordsByUserID(userID)
	if err != nil {
		return err
	}
	return nil
}

// logoutByRefreshToken implements AuthService.
func (a *AuthService) logoutByRefreshToken(refreshToken string) *domain.Error {
	err := a.refreshTokenService.DeleteOneByRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}

// generateAndUpdateToken implements AuthService.
func (a *AuthService) generateAndUpdateToken(user *entity.User, refreshTokenID int) (*auth_dto.NewTokenResponse, *domain.Error) {
	token, tokenErr := a.generateNewJWTAndRefreshToken(user)
	if tokenErr != nil {
		return nil, tokenErr
	}
	now := time.Now()
	update, updateErr := a.refreshTokenService.UpdateOneRefreshTokenByUserID(refreshTokenID, entity.RefreshToken{
		AccessToken: token.AccessToken,
		ExpiredAt:   &token.RefreshTokenExpiredAt,
		Token:       token.RefreshToken,
		UserID:      user.ID,
		UpdatedAt:   &now,
	})
	if updateErr != nil {
		return nil, updateErr
	}
	return &auth_dto.NewTokenResponse{
		AccessToken:           update.AccessToken,
		RefreshToken:          update.Token,
		AccessTokenExpiredAt:  token.AccessTokenExpiredAt.Format(time.RFC3339),
		RefreshTokenExpiredAt: token.RefreshTokenExpiredAt.Format(time.RFC3339),
	}, nil
}

// GenerateNewTokenAndRefreshToken implements AuthService.
func (a *AuthService) GenerateNewTokenAndRefreshToken(refreshToken string) (*auth_dto.NewTokenResponse, *domain.Error) {
	refresh, refreshErr := a.refreshTokenService.CheckSessionValidityByRefreshToken(refreshToken)
	if refreshErr != nil {
		return nil, refreshErr
	}

	user, userErr := a.userService.FindOneByID(refresh.UserID)
	if userErr != nil {
		return nil, userErr
	}

	newToken, newTokenErr := a.generateAndUpdateToken(user, refresh.ID)
	if newTokenErr != nil {
		return nil, newTokenErr
	}
	return newToken, nil
}

// Self implements AuthService.
func (a *AuthService) Self(accessToken string, user *auth_dto.ATClaims) (*auth_dto.SelfResponse, *domain.Error) {
	checkUser, checkUserErr := a.userService.FindOneByID(user.UserID)
	if checkUserErr != nil {
		if errors.Is(checkUserErr.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(fiber.StatusNotFound, errors.New("RECORD_NOT_FOUND"), nil)
		}
		return nil, checkUserErr
	}

	myCircle := (*entity.Circle)(nil)
	if checkUser.CircleID != nil {
		circle, circleErr := a.circleService.GetOneCircleByCircleID(*checkUser.CircleID)

		if circleErr != nil && !errors.Is(circleErr.Err, gorm.ErrRecordNotFound) {
			return nil, circleErr
		}

		if circle != nil {
			myCircle = circle
		}

		if user.CircleID == nil && checkUser.CircleID != nil {
			user.CircleID = checkUser.CircleID
		}
	}

	return &auth_dto.SelfResponse{
		User:                 *checkUser,
		Circle:               myCircle,
		AccessTokenExpiredAt: user.ExpiresAt.Time.Format(time.RFC3339),
	}, nil
}

// GetAuthURL implements AuthService.
func (a *AuthService) GetAuthURL() string {
	return a.config.AuthCodeURL()
}

// registerWithGoogle implements AuthService.
func (a *AuthService) registerWithGoogle(user *auth_dto.GoogleUserData) (*entity.User, *domain.Error) {
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
func (a *AuthService) generateNewJWTAndRefreshToken(user *entity.User) (*auth_dto.NewToken, *domain.Error) {
	appStage := os.Getenv("APP_STAGE")
	secret := os.Getenv("JWT_SECRET")
	var duration time.Duration

	if appStage == "local" {
		duration = time.Minute * 60
	} else {
		duration = time.Minute * 15
	}

	expiredAt := time.Now().Add(duration)
	claims := auth_dto.ATClaims{
		BasicClaims: auth_dto.BasicClaims{
			UserID:   user.ID,
			Email:    user.Email,
			CircleID: user.CircleID,
		},
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
func (a *AuthService) login(user *entity.User) (*auth_dto.NewTokenResponse, *domain.Error) {
	newToken, newTokenErr := a.generateNewJWTAndRefreshToken(user)
	if newTokenErr != nil {
		return nil, newTokenErr
	}

	_, insertErr := a.refreshTokenService.CreateOneRefreshToken(entity.RefreshToken{
		AccessToken: newToken.AccessToken,
		Token:       newToken.RefreshToken,
		UserID:      user.ID,
		ExpiredAt:   &newToken.RefreshTokenExpiredAt,
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
func (a *AuthService) AuthWithGoogle(code string) (*auth_dto.NewTokenResponse, *domain.Error) {
	user, err := a.config.ParseCodeToUserData(code)
	if err != nil {
		return nil, err
	}

	existingUser, existingUserErr := a.userService.FindOneByEmail(user.Email)

	if existingUserErr != nil && !errors.Is(existingUserErr.Err, gorm.ErrRecordNotFound) {
		return nil, existingUserErr
	}

	if existingUser != nil {
		deleteErr := a.refreshTokenService.DeleteAllRefreshTokenRecordsByUserID(existingUser.ID)
		if deleteErr != nil {
			return nil, deleteErr
		}
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
	userService *user.UserService,
	config internal_config.Config,
	refreshToken *refreshtoken.RefreshTokenService,
	utils utils.Utils,
	circleService *circle.CircleService,
) *AuthService {
	return &AuthService{
		userService,
		config,
		refreshToken,
		utils,
		circleService,
	}
}
