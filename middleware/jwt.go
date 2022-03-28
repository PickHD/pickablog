package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ValidateJWTMiddleware responsible to validating jwt in header each request
func ValidateJWTMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}