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
	userService *user.UserService
}

func (a *AuthMiddleware) AdminOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth_dto.ATClaims)
	if user.UserID != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
	}

	return c.Next()
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

func (a *AuthMiddleware) parseToken(accessToken string) (*auth_dto.ATClaims, *domain.Error) {
	secret := os.Getenv("JWT_SECRET")
	claims := &auth_dto.ATClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, domain.NewError(fiber.StatusUnauthorized, errors.New("INVALID_SIGNATURE"), nil)
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.NewError(fiber.StatusUnauthorized, errors.New("TOKEN_EXPIRED"), nil)
		}

		return nil, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)
	}

	if !token.Valid {
		return nil, domain.NewError(fiber.StatusUnauthorized, errors.New("TOKEN_INVALID"), nil)
	}

	return claims, nil
}

func (a *AuthMiddleware) IfAuthed(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	if accessToken == "" {
		c.Locals("user", nil)
		return c.Next()
	}
	claims, err := a.parseToken(accessToken)
	if err != nil {
		type reqCookie struct {
			RefreshToken string `cookie:"refresh_token"`
		}
		reqCookies := new(reqCookie)
		if err := c.CookieParser(reqCookies); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
		}

		if reqCookies.RefreshToken == "" {
			c.Locals("user", nil)
			return c.Next()
		} else {
			return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
		}
	}

	c.Locals("user", claims)
	c.Locals("accessToken", accessToken)
	return c.Next()
}

func (a *AuthMiddleware) Init(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	type reqCookie struct {
		RefreshToken string `cookie:"refresh_token"`
	}
	reqCookies := new(reqCookie)
	if err := c.CookieParser(reqCookies); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	if accessToken == "" {
		if reqCookies.RefreshToken == "" {
			return c.
				Status(fiber.StatusUnauthorized).
				JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
		} else {
			return c.
				Status(fiber.StatusUnauthorized).
				JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("TOKEN_INVALID"), nil)))
		}

	}

	claims, err := a.parseToken(accessToken)
	if err != nil {

		if err.Err.Error() == "TOKEN_INVALID" && reqCookies.RefreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusUnauthorized, errors.New("UNAUTHORIZED"), nil)))
		}

		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	c.Locals("user", claims)
	c.Locals("accessToken", accessToken)

	return c.Next()
}

func NewAuthMiddleware(
	userService *user.UserService,
) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}
