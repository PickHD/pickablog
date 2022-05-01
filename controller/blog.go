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
	// IBlogController is an interface that has all the function to be implemented inside blog controller
	IBlogController interface {
		CreateBlog(ctx *fiber.Ctx) error
		ListBlog(ctx *fiber.Ctx) error
		DetailBlog(ctx *fiber.Ctx) error
		UpdateBlog(ctx *fiber.Ctx) error
		DeleteBlog(ctx *fiber.Ctx) error

		CreateComment(ctx *fiber.Ctx) error
		UpdateComment(ctx *fiber.Ctx) error
		ListComment(ctx *fiber.Ctx) error
		DeleteComment(ctx *fiber.Ctx) error

		Like(ctx *fiber.Ctx) error
		UnLike(ctx *fiber.Ctx) error
	}

	// BlogController is an app blog struct that consists of all the dependencies needed for blog controller
	BlogController struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		BlogSvc service.IBlogService
	}
)
// CreateBlog responsible to creating a blog from controller layer
func (bc *BlogController) CreateBlog(ctx *fiber.Ctx) error {
	var blogReq model.CreateBlogRequest

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&blogReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	blogReq.UserID = extData.UserID

	err = bc.BlogSvc.CreateBlogSvc(blogReq,extData.FullName)
	if err != nil {
		if errors.Is(err,model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusCreated,nil,"Success Create Blog",nil,nil)
}

// ListBlog responsible to listing/filter/search a blogs from controller layer
func (bc *BlogController) ListBlog(ctx *fiber.Ctx) error {
	var (
		page = 1
		size = 10
		order = "ASC"
		field = "id"
		search = ""
		filter  model.FilterBlogRequest
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

	if err := ctx.QueryParser(&filter); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	data,meta,err := bc.BlogSvc.GetAllBlogSvc(page,size,order,field,search,filter)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Getting all Blogs",data,meta)
}

// DetailBlog responsible to getting detail a blog by slug from controller layer
func (bc *BlogController) DetailBlog(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug","")

	if slug == "" {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	data,err := bc.BlogSvc.GetBlogBySlugSvc(slug)
	if err != nil {
		if errors.Is(err,model.ErrBlogNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Get Blog",data,nil)
}

// UpdateBlog responsible to updating a blog by id from controller layer
func (bc *BlogController) UpdateBlog(ctx *fiber.Ctx) error {
	var blogReq model.UpdateBlogRequest

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&blogReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	blogReq.UserID = extData.UserID

	err = bc.BlogSvc.UpdateBlogSvc(id,blogReq,extData.FullName)
	if err != nil {
		if errors.Is(err,model.ErrForbiddenUpdate) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrBlogNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrEmailExisted) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Update Blog",nil,nil)
}

// DeleteBlog responsible to deleting a blog by ud from controller layer
func (bc *BlogController) DeleteBlog(ctx *fiber.Ctx) error {
	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	err = bc.BlogSvc.DeleteBlogSvc(id,extData.UserID)
	if err != nil {
		if errors.Is(err,model.ErrForbiddenDelete) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrBlogNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Delete Blog",nil,nil)
}

// CreateComment responsible to creating a comment from controller layer
func (bc *BlogController) CreateComment(ctx *fiber.Ctx) error {
	var commentReq model.CommentRequest

	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&commentReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	commentReq.UserID = extData.UserID

	err = bc.BlogSvc.CreateCommentSvc(id,commentReq,extData.FullName)
	if err != nil {

		if errors.Is(err,model.ErrUserNotFound) || errors.Is(err,model.ErrBlogNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusCreated,nil,"Success Create Comment",nil,nil)
}

// UpdateComment responsible to updating a comment by id from controller layer
func (bc *BlogController) UpdateComment(ctx *fiber.Ctx) error {
	var commentReq model.CommentRequest

	commentId,err := ctx.ParamsInt("comment_id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if commentId == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrCommentNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	if err := ctx.BodyParser(&commentReq); err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	commentReq.UserID = extData.UserID

	err = bc.BlogSvc.UpdateCommentSvc(commentId,commentReq,extData.FullName)
	if err != nil {
		if errors.Is(err,model.ErrForbiddenUpdate) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrUserNotFound) || errors.Is(err,model.ErrCommentNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Update Comment",nil,nil)
}

// ListComment responsible to getting comments by blog id from controller layer
func (bc *BlogController) ListComment(ctx *fiber.Ctx) error {
	var (
		page = 1
		size = 10
		order = "ASC"
		field = "id"
	) 

	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

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

	data,meta,err := bc.BlogSvc.GetCommentsByBlogSvc(id,page,size,order,field)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Getting all Comments",data,meta)
}

// DeleteComment responsible to deleting a comment by id from controller layer
func (bc *BlogController) DeleteComment(ctx *fiber.Ctx) error {
	blogID,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if blogID == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	commentID,err := ctx.ParamsInt("comment_id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if commentID == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrCommentNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	err = bc.BlogSvc.DeleteCommentSvc(blogID,commentID,extData.UserID)
	if err != nil {
		if errors.Is(err,model.ErrForbiddenDelete) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusForbidden,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrBlogNotFound) || errors.Is(err,model.ErrCommentNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success Delete Comment",nil,nil)
}

// Like responsible to creating a like / liking a blog
func (bc *BlogController) Like(ctx *fiber.Ctx) error {
	var likeReq model.LikeRequest

	id,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if id == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	likeReq.Like = 1
	likeReq.UserID = extData.UserID

	err = bc.BlogSvc.CreateLikeSvc(id,likeReq,extData.FullName)
	if err != nil {

		if errors.Is(err,model.ErrUserNotFound) || errors.Is(err,model.ErrBlogNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		if errors.Is(err,model.ErrInvalidRequest) || errors.Is(err,model.ErrAlreadyLikeBlog) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusCreated,nil,"Success Like",nil,nil)
}

// UnLike responsible to deleting a like by id / unliking a blog
func (bc *BlogController) UnLike(ctx *fiber.Ctx) error {
	blogID,err := ctx.ParamsInt("id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if blogID == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrBlogNotFound.Error(),nil,nil)
	}

	likeID,err := ctx.ParamsInt("like_id",0)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,err,err.Error(),nil,nil)
	}

	if likeID == 0 {
		return helper.ResponseFormatter[any](ctx,fiber.StatusBadRequest,nil,model.ErrLikeNotFound.Error(),nil,nil)
	}

	data := ctx.Locals(model.KeyJWTValidAccess)
	extData,err := util.ExtractPayloadJWT(data)
	if err != nil {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	err = bc.BlogSvc.DeleteLikeSvc(blogID,likeID,extData.UserID)
	if err != nil {
		if errors.Is(err,model.ErrBlogNotFound) || errors.Is(err,model.ErrLikeNotFound) {
			return helper.ResponseFormatter[any](ctx,fiber.StatusNotFound,err,err.Error(),nil,nil)
		}

		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,err.Error(),nil,nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"Success UnLike",nil,nil)
}
