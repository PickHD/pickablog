package controller

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/service"
	"github.com/PickHD/pickablog/util"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	// IUserController is an interface that has all the function to be implemented inside user controller
	IUserController interface {
		ListUser(ctx *fiber.Ctx) error
		DetailUser(ctx *fiber.Ctx) error
		UpdateUser(ctx *fiber.Ctx) error
		DeleteUser(ctx *fiber.Ctx) error
	}
	
	// UserController is an app tag struct that consists of all the dependencies needed for user controller
	UserController struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		UserSvc service.IUserService
	}
)

// ListUser responsible to getting all user from controller layer
func (uc *UserController) ListUser(ctx *fiber.Ctx) error {
	var (
		page = 1
		size = 10
		order = "ASC"
		field = "id"
		search = ""
	)

	if p := ctx.Query("page",""); p != "" {
		pNum,err := strconv.Atoi(p)
		if err == nil && pNum > 0 {
			page = pNum
		}
	}

	if s := ctx.Query("size",""); s != "" {
		sNum, err := strconv.Atoi(s)
		if err == nil && sNum > 0 {
			size = sNum
		}
	}

	if o := ctx.Query("order",""); o != "" {
		if len(o) > 0 {
			order = o
		}
	}

	if f := ctx.Query("field",""); f != "" {
		if len(f) > 0 {
			field = f
		}
	}

	if sr := ctx.Query("s",""); sr != "" {
		if len(sr) > 0 {
			search = strings.Trim(sr," ")
		}
	}

	data,meta,err := uc.UserSvc.GetAllUserSvc(page,size,order,field,search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Getting all Users",data,meta)
}

// DetailUser responsible to get a user by id from controller layer
func (uc *UserController) DetailUser(ctx *fiber.Ctx) error {
	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrUserNotFound.Error(),nil,nil)
	}

	data,err := uc.UserSvc.GetUserByIDSvc(id)
	if err != nil {
		if errors.Is(err,model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Get User",data,nil)
}

// UpdateUser responsible to update a user by id from controller layer
func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var userReq model.UpdateUserRequest

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&userReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrTagNotFound.Error(),nil,nil)
	}

	err = uc.UserSvc.UpdateUserByIDSvc(id,userReq,extData.FullName)
	if err != nil {
		if errors.Is(err,model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrEmailExisted) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Update User", nil, nil)
}

// DeleteUser responsible to delete a user by id from controller layer
func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrUserNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if id == extData.UserID {
		return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,model.ErrForbiddenDeleteSelf,model.ErrForbiddenDeleteSelf.Error(),nil,nil)
	}

	err = uc.UserSvc.DeleteUserByIDSvc(id)
	if err != nil {
		if errors.Is(err,model.ErrUserNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Delete User",nil,nil)
}