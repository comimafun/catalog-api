package auth

import (
	"catalog-be/internal/domain"
	"errors"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService AuthService
	validator   *validator.Validate
}

func (a *AuthHandler) GetAuthURL(c *fiber.Ctx) error {
	url := a.authService.GetAuthURL()
	c.Status(fiber.StatusFound)
	return c.Redirect(url)
}

func (a *AuthHandler) GetGoogleCallback(c *fiber.Ctx) error {
	stateQ := c.Query("state")
	stateEnv := os.Getenv("GOOGLE_STATE")
	if stateQ != stateEnv {
		err := errors.New("INVALID_STATE")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusUnauthorized,
		})
	}
	code := c.Query("code")
	if code == "" {
		err := errors.New("INVALID_CODE")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	data, err := a.authService.AuthWithGoogle(code)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})
}

func (a *AuthHandler) PostGoogleCallback(c *fiber.Ctx) error {
	type reqBody struct {
		Code string `json:"code" validate:"required"`
	}

	code := new(reqBody)
	if err := c.BodyParser(code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	if err := a.validator.Struct(code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.NewErrorFiber(c, domain.NewError(fiber.StatusBadRequest, err, nil)))
	}

	data, err := a.authService.AuthWithGoogle(code.Code)
	if err != nil {
		return c.Status(err.Code).JSON(domain.NewErrorFiber(c, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})

}

func NewAuthHandler(authService AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService,
		validator,
	}
}
