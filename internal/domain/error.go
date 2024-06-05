package domain

import "github.com/gofiber/fiber/v2"

type Error struct {
	Code int     `json:"code"`
	Err  error   `json:"error"`
	Info *string `json:"info"`
}

func NewError(code int, err error, info *string) *Error {
	return &Error{
		Code: code,
		Err:  err,
		Info: info,
	}
}

func NewErrorFiber(ctx *fiber.Ctx, err *Error) fiber.Map {
	return fiber.Map{
		"error_id": ctx.Locals("requestid"),
		"error":    err.Err.Error(),
		"code":     err.Code,
	}
}
