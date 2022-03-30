package middleware

import (
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/util"
	"github.com/gofiber/fiber/v2"
)

// SuperAdminOnlyMiddleware responsible to ensure authorized role is superadmin
func SuperAdminOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	if extData.RoleName == "Superadmin" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,model.ErrForbiddenAccess,model.ErrForbiddenAccess.Error(),nil)
}

// SuperAdminOnlyMiddleware responsible to ensure authorized role is author
func AuthorOnlyMiddleware(ctx *fiber.Ctx) error {
	data := ctx.Locals(model.KeyJWTValidAccess)

	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	if extData.RoleName == "Author" {
		return ctx.Next()
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,model.ErrForbiddenAccess,model.ErrForbiddenAccess.Error(),nil)
}