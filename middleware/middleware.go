package middleware

import "github.com/gofiber/fiber/v2"

type middleware func(fiber.Handler) fiber.Handler

// Chain function purpose to chain multiple handler middleware
func Chain(h fiber.Handler,m ...middleware) fiber.Handler {
	return applyChaining(h,m...)
}

// initialize middleware for function
func applyChaining(h fiber.Handler,m ...middleware) fiber.Handler {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	for i := len(m) - 1;i >=0; i-- {
		wrapped = m[i](wrapped)
	} 

	return wrapped
}