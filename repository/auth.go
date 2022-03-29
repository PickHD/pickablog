package repository

import (
	"context"
	"fmt"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/model"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	IAuthRepository interface {
		CreateUser(user model.RegisterAuthorRequest) error
		GetUserByEmail(email string) (*model.AuthUserDetails,error)
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// CreateUser repository layer for executing command creating a user
func (ar *AuthRepository) CreateUser(user model.RegisterAuthorRequest) error {
	q := `INSERT INTO "user" (full_name,email,password,role_id,created_by) VALUES ($1,$2,$3,$4,$5)`

	_,err := ar.DB.Exec(ar.Context,q,user.FullName,user.Email,user.Password,2,user.FullName)
	if err != nil {
		ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser ERROR : %v",err))
		return err
	}

	return nil
}

// GetUserByEmail repository layer for querying command getting any user by email
func (ar *AuthRepository) GetUserByEmail(email string) (*model.AuthUserDetails,error) {
	var authUserDetail model.AuthUserDetails

	q := `SELECT 
		u.full_name,
		u.email,
		r.id,
		r.name FROM "user" u 
		LEFT JOIN role r ON r.id = u.role_id 
		WHERE u.email = $1
	`

	row := ar.DB.QueryRow(ar.Context,q,email)
	err := row.Scan(&authUserDetail.FullName,&authUserDetail.Email,&authUserDetail.RoleID,&authUserDetail.RoleName)
	if err != nil {
		ar.Logger.Error(fmt.Errorf("AuthRepository.GetUserByEmail ERROR : %v",err))
		return nil,err
	}

	return &authUserDetail,nil
}