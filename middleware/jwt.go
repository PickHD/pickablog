package middleware

import (
	"fmt"

	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/util"
	"github.com/gofiber/fiber/v2"
)

// ValidateJWTMiddleware responsible to validating jwt in header each request
func ValidateJWTMiddleware(ctx *fiber.Ctx) error {
	// validate JWT coming from request, if valid decode into a struct
	decodedPayload,err := util.ValidateJWT(ctx)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusUnauthorized,err, fmt.Sprintf("Unauthorized access, reason : %s",err.Error()),nil,nil)
	}

	// pass decoded payload into ctx.Locals()
	ctx.Locals(model.KeyJWTValidAccess,decodedPayload)

	// going to next handler..
	return ctx.Next()
}