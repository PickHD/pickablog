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
		GoogleLogin(ctx *fiber.Ctx) error
		GoogleLoginCallback(ctx *fiber.Ctx) error
		Login(ctx *fiber.Ctx) error
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
	var regAuthorReq model.CreateUserRequest

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

// GoogleLogin responsible to handling redirect to google auth services from controller layer
func (ac *AuthController) GoogleLogin(ctx *fiber.Ctx) error {
	url,err := ac.AuthSvc.GoogleLoginSvc()
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	return ctx.Redirect(url)
}

// GoogleLoginCallback responsible to handling callback redirect to API from google auth services from controller layer
func (ac *AuthController) GoogleLoginCallback(ctx *fiber.Ctx) error {
	getState := ctx.Query("state","")
	getCode := ctx.Query("code","")

	if len(getState) <= 0 || len(getCode) <= 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusNotAcceptable,model.ErrInvalidExchange,model.ErrInvalidExchange.Error(),nil)
	}

	jwt,err := ac.AuthSvc.GoogleLoginCallbackSvc(getState,getCode)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Login",jwt)
}

// Login respoinsible to log-in data from controller layer
func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	var loginReq model.LoginRequest

	if err := ctx.BodyParser(&loginReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,model.ErrFailedParseBody.Error(),nil)
	}

	jwt,err := ac.AuthSvc.LoginSvc(loginReq)
	if err != nil {
		if errors.Is(err,model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrInvalidPassword) || errors.Is(err,model.ErrMismatchLogin) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Login",jwt)
}