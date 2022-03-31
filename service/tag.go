package service

import (
	"context"

	"github.com/PickHD/pickablog/config"
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
func (tr *TagService) CreateTagSvc(req model.CreateTagRequest,createdBy string) error {
	err := validateCreateTagRequest(&req)
	if err != nil {
		return err
	}

	_,err = tr.TagRepo.GetByName(req.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = tr.TagRepo.Create(req,createdBy)
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
func (tr *TagService) GetAllTagSvc(page int,size int,order string, field string, search string) ([]model.ViewTagResponse,*model.Metadata,error) {
	data,totalData, err := tr.TagRepo.GetAll(page,size,order,field,search)
	if err != nil {
		return nil,nil,err
	}

	totalPage := (int(totalData) + size - 1) / size

	if len(data) > size {
		data = data[:len(data)-1]
	}

	meta := buildTagMetaData(page,size,order,totalData,totalPage)

	return data,meta,nil
} 

// validateCreateTagRequest responsible to validating request create tag
func validateCreateTagRequest(req *model.CreateTagRequest) error {
	if len(req.Name) < 5 {
		return model.ErrInvalidRequest
	}

	return nil
}

// buildTagMetaData responsible to building response meta get all tag
func buildTagMetaData(page int,size int,order string,totalData int,totalPage int) *model.Metadata {
	return &model.Metadata{
		Page: page,
		Size: size,
		Order: order,
		TotalData: totalData,
		TotalPage: totalPage,
	}
}