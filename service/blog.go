package service

import (
	"context"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/repository"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IBlogService is an interface that has all the function to be implemented inside blog service
	IBlogService interface {
		CreateBlogSvc(req model.CreateBlogRequest, createdBy string) error
		GetAllBlogSvc(page int, size int, order string, field string, search string, filter model.FilterBlogRequest) ([]model.ViewBlogResponse,*model.Metadata,error)
		GetBlogBySlugSvc(slug string) (*model.ViewBlogResponse,error)
		UpdateBlogSvc(id int, req model.UpdateBlogRequest, updatedBy string) error
		DeleteBlogSvc(id int) error
	}

	// BlogRepository is an app blog struct that consists of all the dependencies needed for blog service
	BlogService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		BlogRepo repository.IBlogRepository
	}
)

// CreateBlogSvc service layer for handling creating a blog
func (bs *BlogService) CreateBlogSvc(req model.CreateBlogRequest, createdBy string) error {
	err := validateCreateBlogRequest(&req)
	if err != nil {
		return err
	}

	req.Slug = helper.GenerateSlug(req.Title)

	err = bs.BlogRepo.Create(req,createdBy)
	if err != nil {
		return err
	}

	return nil
}

// GetAllBlogSvc service layer for handling list/filter/search a blog
func (bs *BlogService) GetAllBlogSvc(page int, size int, order string, field string, search string,filter model.FilterBlogRequest) ([]model.ViewBlogResponse,*model.Metadata,error) {
	data,totalData,err := bs.BlogRepo.GetAll(page,size,order,field,search,filter)
	if err != nil {
		return nil,nil,err
	}

	if totalData < 1 {
		return []model.ViewBlogResponse{},nil,nil
	}

	totalPage := (int(totalData) + size - 1) / size

	if len(data) > size {
		data = data[:len(data)-1]
	}

	meta := helper.BuildMetaData(page,size,order,totalData,totalPage)

	return data,meta,nil
}

// GetBlogBySlugSvc service layer for handling getting detail a blog by slugs
func (bs *BlogService) GetBlogBySlugSvc(slug string) (*model.ViewBlogResponse,error) {
	data,err := bs.BlogRepo.GetBySlug(slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil,model.ErrBlogNotFound
		}

		return nil,err
	}

	return data,nil
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

	req.Slug = helper.GenerateSlug(req.Title)

	blogMap["title"] = req.Title
	blogMap["slug"] = req.Slug
	blogMap["body"] = req.Body
	blogMap["footer"] = req.Footer
	blogMap["user_id"] = req.UserID

	return blogMap,nil
}