package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/PickHD/pickablog/repository"
	"github.com/PickHD/pickablog/requester"
	"github.com/PickHD/pickablog/util"
	"github.com/jackc/pgx/v4"
	"golang.org/x/oauth2"

	"github.com/sirupsen/logrus"
)

type (
	// IAuthService is an interface that has all the function to be implemented inside auth service
	IAuthService interface {
		Create(user model.CreateUserRequest) error
		GoogleLoginSvc() (string,error)
		GoogleLoginCallbackSvc(state string, code string) (*model.SuccessLoginResponse,error)
		LoginSvc(user model.LoginRequest) (*model.SuccessLoginResponse,error)
	}

	// AuthService is an app auth struct that consists of all the dependencies needed for auth service
	AuthService struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		AuthRepo repository.IAuthRepository
		GConfig *oauth2.Config
		GOAuthReq requester.IOAuthGoogle
	}
)

// Create service layer for handling create a user
func (as *AuthService) Create(user model.CreateUserRequest) error {
	err := validateRegisterAuthorRequest(&user)
	if err != nil {
		return err
	}

	_,err = as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			user.Password,err = helper.HashPassword(user.Password)
			if err != nil {
				return err
			}

			err = as.AuthRepo.CreateUser(user,model.RoleAuthor)
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

// GoogleLoginSvc service layer for handling googles login services
func (as *AuthService) GoogleLoginSvc() (string,error) {
	state := helper.GenerateRandomToken(8)

	err := as.AuthRepo.SetRedis(fmt.Sprintf("%s:%s",model.OauthStateKey,state),state, time.Minute * time.Duration(as.Config.Redis.RDBExpire))
	if err != nil {
		return "",err
	}

	url := as.GConfig.AuthCodeURL(state)

	return url, nil
}

// GoogleLoginCallbackSvc service layer for handling googles login services callback
func (as *AuthService) GoogleLoginCallbackSvc(state string,code string) (*model.SuccessLoginResponse,error) {
	// ensure state is valid & stil exists in redis
	_,err := as.AuthRepo.GetRedis(fmt.Sprintf("%s:%s",model.OauthStateKey,state))
	if err != nil {
		return nil,err
	}

	// gain user info based on code from requester oauth google
	gOauthUser,err := as.GOAuthReq.GetUserInfo(code)
	if err != nil {
		return nil,err
	}

	getUser,err := as.AuthRepo.GetUserByEmail(gOauthUser.Email)
	if err != nil {
		// if didnt exists, insert this user to database, then generate JWT
		if err == pgx.ErrNoRows {

			newUser := model.CreateUserRequest{
				FullName:gOauthUser.Name,
				Email:gOauthUser.Email,
				Password: "",
			}

			err = as.AuthRepo.CreateUser(newUser,model.RoleGuest)
			if err != nil {
				return nil,err
			}

			getRoleName,err := model.GetValidRoleByID(model.RoleGuest)
			if err != nil {
				return nil,err
			}

			jwt,err := util.BuildJWT(as.Config,&model.AuthUserDetails{FullName: newUser.FullName,Email:newUser.Email,RoleName: getRoleName})
			if err != nil {
				as.Logger.Error(fmt.Errorf("AuthService.BuildJWT ERROR : %v MSG : %s",err,err.Error()))
				return nil,err
			}

			return &model.SuccessLoginResponse{
				AccessToken: jwt,
				ExpiredAt: time.Now().Add(util.JWTExpire),
				Role: getRoleName,
			},nil
		} 
		
		return nil,err
	}

	// if exists then regenerate JWT
	jwt,err := util.BuildJWT(as.Config,getUser)
	if err != nil {
		as.Logger.Error(fmt.Errorf("AuthService.BuildJWT ERROR : %v MSG : %s",err,err.Error()))
		return nil,err
	}

	return &model.SuccessLoginResponse{
		AccessToken: jwt,
		ExpiredAt: time.Now().Add(util.JWTExpire),
		Role: getUser.RoleName,
	},nil
}

// LoginSvc service layer for handling user login (author,superadmin)
func (as *AuthService) LoginSvc(user model.LoginRequest) (*model.SuccessLoginResponse,error) {
	err := validateLoginRequest(&user)
	if err != nil {
		return nil,err
	}

	getUser,err := as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil,model.ErrUserNotFound
		}

		return nil,err
	}

	// when password is null, this user is using OAUTH login method
	if getUser.Password == "" {
		return nil,model.ErrMismatchLogin
	}

	if !helper.CheckPasswordHash(getUser.Password,user.Password) {
		return nil,model.ErrInvalidPassword
	}

	jwt,err := util.BuildJWT(as.Config,getUser)
	if err != nil {
		as.Logger.Error(fmt.Errorf("AuthService.BuildJWT ERROR : %v MSG : %s",err,err.Error()))
		return nil,err
	}
	
	return &model.SuccessLoginResponse{
		AccessToken: jwt,
		ExpiredAt: time.Now().Add(util.JWTExpire),
		Role: getUser.RoleName,
	},nil
}

// validateRegisterAuthorRequest responsible to validating request register author
func validateRegisterAuthorRequest(user *model.CreateUserRequest) error {
	if len(user.FullName) < 5 || len(user.Password) < 5{
		return model.ErrInvalidRequest
	}

	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return model.ErrInvalidRequest
	}

	return nil
}

// validateLoginRequest responsible to validating request login data
func validateLoginRequest(user *model.LoginRequest) error {
	if !model.IsAllowedEmailInput.MatchString(user.Email) {
		return model.ErrInvalidRequest
	}

	return nil
}