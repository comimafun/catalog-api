package middlewares

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	"catalog-be/internal/modules/user"
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	userService user.UserService
}

func (a *AuthMiddleware) CircleOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)
	if user.CircleID == nil {
		found, err := a.userService.FindOneByID(user.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.NewErrorFiber(c, err))
		}

		if found.CircleID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
		}

		user.CircleID = found.CircleID
		c.Locals("user", user)
	}

	return c.Next()
}

func (a *AuthMiddleware) Init(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
	}
	secret := os.Getenv("JWT_SECRET")
	claims := &auth_dto.ATClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("INVALID_SIGNATURE"), nil)))
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("TOKEN_EXPIRED"), nil)))
		}

		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("TOKEN_INVALID"), nil)))
	}

	c.Locals("user", claims)
	c.Locals("accessToken", accessToken)

	return c.Next()
}

func NewAuthMiddleware(
	userService user.UserService,
) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}
