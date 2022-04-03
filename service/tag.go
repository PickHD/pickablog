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
	// IHealthCheckService is an interface that has all the function to be implemented inside tag service
	ITagService interface {
		CreateTagSvc(req model.CreateTagRequest,createdBy string) error
		GetAllTagSvc(page int,size int,order string,field string, search string) ([]model.ViewTagResponse,*model.Metadata,error)
		UpdateTagSvc(id int, req model.UpdateTagRequest, updatedBy string) error
		DeleteTagSvc(id int) error
	}

	// TagService is an app tag struct that consists of all the dependencies needed for tag service
	TagService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		TagRepo repository.ITagRepository
	}
)

// CreateTagSvc service layer for creating a tag
func (ts *TagService) CreateTagSvc(req model.CreateTagRequest,createdBy string) error {
	err := validateCreateTagRequest(&req)
	if err != nil {
		return err
	}

	_,err = ts.TagRepo.GetByName(req.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = ts.TagRepo.Create(req,createdBy)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	return model.ErrTagNameExisted
}

// GetAllTagSvc service layer for getting all tag
func (ts *TagService) GetAllTagSvc(page int,size int,order string, field string, search string) ([]model.ViewTagResponse,*model.Metadata,error) {
	data,totalData, err := ts.TagRepo.GetAll(page,size,order,field,search)
	if err != nil {
		return nil,nil,err
	}

	if totalData < 1 {
		return []model.ViewTagResponse{},nil,nil
	}

	totalPage := (int(totalData) + size - 1) / size

	if len(data) > size {
		data = data[:len(data)-1]
	}

	meta := helper.BuildMetaData(page,size,order,totalData,totalPage)

	return data,meta,nil
} 

// UpdateTagSvc service layer for updating a tag by id
func (ts *TagService) UpdateTagSvc(id int, req model.UpdateTagRequest, updatedBy string) error {
	tag,err := ts.TagRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrTagNotFound
		}

		return err
	}

	tagMap,err := validateUpdateTagRequest(tag,&req)
	if err != nil {
		return err
	}

	err = ts.TagRepo.UpdateByID(id,tagMap,updatedBy)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTagSvc service layer for deleting a tag by id
func (ts *TagService) DeleteTagSvc(id int) error {
	_,err := ts.TagRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrTagNotFound
		}

		return err
	}

	err = ts.TagRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// validateCreateTagRequest responsible to validating request create tag
func validateCreateTagRequest(req *model.CreateTagRequest) error {
	if len(req.Name) < 5 {
		return model.ErrInvalidRequest
	}

	return nil
}

// validateUpdateTagRequest reposible to validating request update tag
func validateUpdateTagRequest(tag *model.ViewTagResponse,req *model.UpdateTagRequest) (map[string]interface{},error) {
	tagMap:= make(map[string]interface{})

	if req.Name != "" {
		if len(req.Name) < 5 {
			return nil,model.ErrInvalidRequest
		}
		
		tagMap["name"] = req.Name
	}

	return tagMap,nil
}