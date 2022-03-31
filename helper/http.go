package helper

import (
	"github.com/PickHD/pickablog/model"
	"github.com/gofiber/fiber/v2"
)

// ResponseFormatter returning formatted JSON responses
func ResponseFormatter [T any] (ctx *fiber.Ctx,statusCode int, err error, message string, data T,meta *model.Metadata) error {
	ctx.Accepts("application/json")

	if statusCode < 400 {
		return ctx.Status(statusCode).JSON(&fiber.Map{
			"status_code": statusCode,
			"message":message,
			"error":nil,
			"data": data,
			"meta": meta,
		})
		
	}

	return ctx.Status(statusCode).JSON(&fiber.Map{
		"status_code": statusCode,
		"message":message,
		"error": err,
		"data": nil,
		"meta": nil,
	})
}

// OptionsHandler will handing preflight requests
func OptionsHandler (ctx *fiber.Ctx) error {return nil}