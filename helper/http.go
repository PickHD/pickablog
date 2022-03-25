package helper

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseFormatter returning formatted JSON responses
func ResponseFormatter [T any] (ctx *fiber.Ctx,statusCode int, err error, message string, data T) error {
	ctx.Accepts("application/json")

	if statusCode < 400 {
		return ctx.JSON(&fiber.Map{
			"status_code": statusCode,
			"message":message,
			"error":nil,
			"data": data,
		})
		
	}

	return ctx.JSON(&fiber.Map{
		"status_code": statusCode,
		"message":message,
		"error": err,
		"data": nil,
	})
}

// OptionsHandler will handing preflight 
func OptionsHandler (ctx *fiber.Ctx) error {return nil}