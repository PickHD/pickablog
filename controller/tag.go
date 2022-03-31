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
	// ITagController is an interface that has all the function to be implemented inside tag controller
	ITagController interface {
		CreateTag(ctx *fiber.Ctx) error
		GetAllTag(ctx *fiber.Ctx) error
	}
	
	// TagController is an app tag struct that consists of all the dependencies needed for tag controller
	TagController struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		TagSvc service.ITagService
	}

)

// CreateTag responsible to creating a tag of article from controller layer/
func (tc *TagController) CreateTag(ctx *fiber.Ctx) error {
	var tagReq model.CreateTagRequest

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&tagReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	err = tc.TagSvc.CreateTagSvc(tagReq,extData.FullName)
	if err != nil {
		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrTagNameExisted) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}
		
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusCreated,nil,"Success Create Tag",nil,nil)
}

// GetAllTag responsible to getting all tag of article from controller layer
func (tc *TagController) GetAllTag(ctx *fiber.Ctx) error {
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

	if s := ctx.Query("s",""); s != "" {
		if len(s) > 0 {
			search = strings.Trim(s," ")
		}
	}

	data,meta,err := tc.TagSvc.GetAllTagSvc(page,size,order,field,search)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Getting all Tags",data,meta)
}