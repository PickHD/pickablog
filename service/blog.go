package service

import (
	"context"
	"sync"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/repository"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IBlogService is an interface that has all the function to be implemented inside blog service
	IBlogService interface {
		CreateBlogSvc(req model.CreateBlogRequest, createdBy string) error
		GetAllBlogSvc(page int, size int, order string, field string, search string, filter model.FilterBlogRequest) ([]model.BlogResponse,*model.Metadata,error)
		GetBlogBySlugSvc(slug string) (*model.BlogResponse,error)
		UpdateBlogSvc(id int, req model.UpdateBlogRequest, updatedBy string) error
		DeleteBlogSvc(id int) error

		CreateCommentSvc(blogID int,req model.CommentRequest,createdBy string) error
		UpdateCommentSvc(id int, req model.CommentRequest, updatedBy string) error
		GetCommentsByBlogSvc(blogID int,page int, size int, order string, field string) ([]model.ViewCommentResponse,*model.Metadata,error)
		DeleteCommentSvc(blogID int, commentID int) error

		CreateLikeSvc(blogID int,req model.LikeRequest,createdBy string) error
		DeleteLikeSvc(blogID int,likeID int) error
	}

	// BlogRepository is an app blog struct that consists of all the dependencies needed for blog service
	BlogService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		BlogRepo repository.IBlogRepository
		CommentRepo repository.ICommentRepository
		LikeRepo repository.ILikeRepository
		Mutex *sync.RWMutex
	}
)

// CreateBlogSvc service layer for handling creating a blog
func (bs *BlogService) CreateBlogSvc(req model.CreateBlogRequest, createdBy string) error {
	err := validateCreateBlogRequest(&req)
	if err != nil {
		return err
	}

	req.Slug = slug.Make(req.Title)

	err = bs.BlogRepo.Create(req,createdBy)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBlogSvc service layer for handling list/filter/search a blog
func (bs *BlogService) GetAllBlogSvc(page int, size int, order string, field string, search string,filter model.FilterBlogRequest) ([]model.BlogResponse,*model.Metadata,error) {
	res,totalData,err := bs.BlogRepo.GetAll(page,size,order,field,search,filter)
	if err != nil {
		return nil,nil,err
	}

	if totalData < 1 {
		return []model.BlogResponse{},nil,nil
	}

	totalPage := (int(totalData) + size - 1) / size

	if len(res) > size {
		res = res[:len(res)-1]
	}

	meta := helper.BuildMetaData(page,size,order,totalData,totalPage)

	var data []model.BlogResponse

	for i := range res {
		r := res[i]
		d := model.BlogResponse{}

		d.ID= r.ID
		d.Title = r.Title
		d.Slug = r.Slug
		d.Body = r.Body
		d.Footer = r.Footer
		d.UserID = r.UserID
		d.CreatedAt = r.CreatedAt
		d.CreatedBy = r.CreatedBy
		d.UpdatedAt = r.UpdatedAt
		d.UpdatedBy = r.UpdatedBy

		if len(r.Comments) > 1 {
			for c := range r.Comments{
				rc := r.Comments[c]
				if rc.Valid {
					d.Comments = append(d.Comments, int(rc.Int32))
				}
			}
		}
		

		if len(r.Tags) > 1 {
			for t := range r.Tags{
				rt := r.Tags[t]
				if rt.Valid {
					d.Tags = append(d.Tags, int(rt.Int32))
				}
			}
		}

		if len(r.Likes) > 1 {
			for l := range r.Likes{
				rl := r.Likes[l]
				if rl.Valid {
					d.Likes = append(d.Likes, int(rl.Int32))
				}
			}
		}

		data = append(data, d)
	}

	return data,meta,nil
}

// GetBlogBySlugSvc service layer for handling getting detail a blog by slugs
func (bs *BlogService) GetBlogBySlugSvc(slug string) (*model.BlogResponse,error) {
	res,err := bs.BlogRepo.GetBySlug(slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil,model.ErrBlogNotFound
		}

		return nil,err
	}

	r := res
	data := model.BlogResponse{}

	data.ID= r.ID
	data.Title = r.Title
	data.Slug = r.Slug
	data.Body = r.Body
	data.Footer = r.Footer
	data.UserID = r.UserID
	data.CreatedAt = r.CreatedAt
	data.CreatedBy = r.CreatedBy
	data.UpdatedAt = r.UpdatedAt
	data.UpdatedBy = r.UpdatedBy

	if len(r.Comments) > 1 {
		for c := range r.Comments{
			rc := r.Comments[c]
			if rc.Valid {
				data.Comments = append(data.Comments, int(rc.Int32))
			}
		}
	}
	

	if len(r.Tags) > 1 {
		for t := range r.Tags{
			rt := r.Tags[t]
			if rt.Valid {
				data.Tags = append(data.Tags, int(rt.Int32))
			}
		}
	}

	if len(r.Likes) > 1 {
		for l := range r.Likes{
			rl := r.Likes[l]
			if rl.Valid {
				data.Likes = append(data.Likes, int(rl.Int32))
			}
		}
	}

	return &data,nil
}

// UpdateBlogSvc service layer for handling updating a blog by ID
func (bs *BlogService) UpdateBlogSvc(id int, req model.UpdateBlogRequest, updatedBy string) error {
	_,err := bs.BlogRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrBlogNotFound
		}
	}

	blogMap,err := validateUpdateBlogRequest(&req)
	if err != nil {
		return err
	}

	err = bs.BlogRepo.UpdateByID(id,blogMap,updatedBy)
	if err != nil {
		return err
	}

	return nil
}

// DeleteBlogSvc service layer for handling deleting a blog by ID
func (bs *BlogService) DeleteBlogSvc(id int) error {
	_,err := bs.BlogRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrBlogNotFound
		}
	}

	err = bs.BlogRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// CreateCommentSvc service layer for handling creating comments
func (bs *BlogService) CreateCommentSvc(blogID int, req model.CommentRequest,createdBy string) error {return nil}

// UpdateCommentSvc service layer for handling updating comment by id
func (bs *BlogService) UpdateCommentSvc(id int, req model.CommentRequest,updatedBy string) error {return nil}

// GetCommentsByBlogSvc service layer for handling getting comments with filter
func (bs *BlogService) GetCommentsByBlogSvc(blogID int, page int, size int, order string, field string) ([]model.ViewCommentResponse,*model.Metadata,error) {return nil,nil,nil}

// DeleteCommentSvc service layer for handling deleting comment with id
func (bs *BlogService) DeleteCommentSvc(blogID int, commentID int) error {return nil}

// CreateLikeSvc service layer for handling creating likes
func (bs *BlogService) CreateLikeSvc(blogID int, req model.LikeRequest,createdBy string) error {return nil}

// DeleteLikeSvc service layer for handling deleting like with id
func (bs *BlogService) DeleteLikeSvc(blogID int,likeID int) error {return nil}

// validateCreateBlogRequest responsible to validating create blog request
func validateCreateBlogRequest(req *model.CreateBlogRequest) error {
	if req.Title == "" || req.Body == "" || req.Footer == "" || len(req.Tags) < 1 {
		return model.ErrInvalidRequest
	}

	if len(req.Title) < 5 && len(req.Title) > 50 {
		return model.ErrInvalidRequest
	}

	if len(req.Body) < 100 {
		return model.ErrInvalidRequest
	}

	if len(req.Footer) < 5 {
		return model.ErrInvalidRequest
	}

	return nil
}

// validateCreateBlogRequest responsible to validating update blog request
func validateUpdateBlogRequest(req *model.UpdateBlogRequest) (map[string]interface{},error) {
	blogMap := make(map[string]interface{})

	if req.Title == "" || req.Body == "" || req.Footer == "" {
		return nil,model.ErrInvalidRequest
	}

	if len(req.Title) < 5 && len(req.Title) > 50 {
		return nil,model.ErrInvalidRequest
	}

	if len(req.Body) < 100 {
		return nil,model.ErrInvalidRequest
	}

	if len(req.Footer) < 5 {
		return nil,model.ErrInvalidRequest
	}

	req.Slug = slug.Make(req.Title)

	blogMap["title"] = req.Title
	blogMap["slug"] = req.Slug
	blogMap["body"] = req.Body
	blogMap["footer"] = req.Footer
	blogMap["user_id"] = req.UserID

	return blogMap,nil
}