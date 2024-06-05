package middlewares

import (
	"catalog-be/internal/domain"
	auth_dto "catalog-be/internal/modules/auth/dto"
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) CircleOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)
	if user.CircleID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
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

	return c.Next()
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}
