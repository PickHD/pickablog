package controller

import (
	"context"
	"errors"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	// IAuthController is an interface that has all the function to be implemented inside auth controller
	IAuthController interface {
		RegisterAuthor(ctx *fiber.Ctx) error
	}

	// AuthController is an app auth struct that consists of all the dependencies needed for auth controller
	AuthController struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		AuthSvc service.IAuthService
	}
)

// RegisterAuthor responsible to registering data author from controller layer
func (ac *AuthController) RegisterAuthor(ctx *fiber.Ctx) error {
	var regAuthorReq model.RegisterAuthorRequest

	if err := ctx.BodyParser(&regAuthorReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,model.ErrFailedParseBody.Error(),nil)
	}
	
	err := ac.AuthSvc.Create(regAuthorReq)
	if err != nil {
		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrEmailExisted) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil)
		}
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusCreated,nil,"Successfully register as author",nil)
}