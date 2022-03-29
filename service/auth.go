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
	// IAuthService is an interface that has all the function to be implemented inside auth service
	IAuthService interface {
		Create(user model.RegisterAuthorRequest) error
	}

	// AuthService is an app auth struct that consists of all the dependencies needed for auth service
	AuthService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		AuthRepo repository.IAuthRepository
	}
)

// Create service layer for handling create a user
func (as *AuthService) Create(user model.RegisterAuthorRequest) error {
	err := validateRegisterAuthorRequest(&user)
	if err != nil {
		return err
	}

	user.Password,err = helper.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_,err = as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = as.AuthRepo.CreateUser(user)
			if err != nil {
				return err
			}

			//TODO : DO WE NEED EMAIL VERIF? using gomail v2
			
			return nil
		}
		return err
	}

	return model.ErrEmailExisted 
}

// validateRegisterAuthorRequest responsible to validating request register author
func validateRegisterAuthorRequest(user *model.RegisterAuthorRequest) error {
	if len(user.FullName) < 5 || len(user.Password) < 5{
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return model.ErrInvalidRequest
	}

	return nil
}

