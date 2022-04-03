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
	// IUserService is an interface that has all the function to be implemented inside user service
	IUserService interface {
		GetAllUserSvc(page int,size int,order string,field string, search string) ([]model.ViewUserResponse,*model.Metadata,error)
		GetUserByIDSvc(id int) (*model.ViewUserResponse,error)
		UpdateUserByIDSvc(id int, req model.UpdateUserRequest,updatedBy string) error
		DeleteUserByIDSvc(id int) error
	}

	// UserService is an app user check struct that consists of all the dependencies needed for user service
	UserService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		UserRepo repository.IUserRepository
	}
)

// GetAllUserSvc service layer for getting all user
func (us *UserService) GetAllUserSvc(page int, size int, order string, field string, search string) ([]model.ViewUserResponse,*model.Metadata,error) {
	data,totalData,err := us.UserRepo.GetAll(page,size,order,field,search)
	if err != nil {
		return nil,nil,err
	}

	if totalData < 1 {
		return []model.ViewUserResponse{},nil,nil
	}

	totalPage := (int(totalData) + size - 1) / size

	if len(data) > size {
		data = data[:len(data)-1]
	}

	meta := helper.BuildMetaData(page,size,order,totalData,totalPage)

	return data,meta,nil
}

// GetUserByIDSvc service layer for get a user by id
func (us *UserService) GetUserByIDSvc(id int) (*model.ViewUserResponse,error) {
	data,err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil,model.ErrUserNotFound
		}

		return nil,err
	}

	return data,nil
}

// UpdateUserByIDSvc service layer for update user by id
func (us *UserService) UpdateUserByIDSvc(id int, req model.UpdateUserRequest,updatedBy string) error {
	data,err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrUserNotFound
		}

		return err
	}

	_,err = us.UserRepo.GetByEmail(req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			userMap,err := validateUpdateUserRequest(data,&req)
			if err != nil {
				return err
			}

			err = us.UserRepo.UpdateByID(id,userMap,updatedBy)
			if err != nil {
				return err
			}

			return nil

		}

		return err
	}

	return model.ErrEmailExisted
}

// DeleteUserByIDSvc service layer for delete user by id
func (us *UserService) DeleteUserByIDSvc(id int) error {
	_,err := us.UserRepo.GetByID(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.ErrUserNotFound
		}

		return err
	}

	err = us.UserRepo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// validateUpdateUserRequest responsible to validating update user
func validateUpdateUserRequest(user *model.ViewUserResponse,req *model.UpdateUserRequest) (map[string]interface{},error) {
	userMap := make(map[string]interface{})

	if len(user.FullName) < 5 {
		return nil,model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return nil,model.ErrInvalidRequest
	}

	userMap["full_name"] = req.FullName
	userMap["email"] = req.Email

	return userMap,nil
}