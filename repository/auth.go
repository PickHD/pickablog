package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/model"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	IAuthRepository interface {
		CreateUser(user model.CreateUserRequest, roleID int) error
		GetUserByEmail(email string) (*model.AuthUserDetails,error)
		SetRedis(key string,value string,expr time.Duration) error
		GetRedis(key string) (string,error)
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
		Redis *redis.Client
	}
)

// CreateUser repository layer for executing command creating a user
func (ar *AuthRepository) CreateUser(user model.CreateUserRequest,roleID int) error {
	q := `INSERT INTO "user" (full_name,email,password,role_id,created_by) VALUES ($1,$2,$3,$4,$5)`

	_,err := ar.DB.Exec(ar.Context,q,user.FullName,user.Email,user.Password,roleID,user.FullName)
	if err != nil {
		ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser ERROR : %v MSG : %s",err,err.Error()))
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
		if err == pgx.ErrNoRows {
			ar.Logger.Info(fmt.Errorf("AuthRepository.GetUserByEmail INFO : %v MSG : %s",err,err.Error()))
		} else {
			ar.Logger.Error(fmt.Errorf("AuthRepository.GetUserByEmail ERROR : %v MSG : %s",err,err.Error()))
		}

		return nil,err
	}

	return &authUserDetail,nil
}

// SetRedis repository layer for set a value into redis by unique key
func (ar *AuthRepository) SetRedis(key string, value string, expr time.Duration) error {
	err := ar.Redis.SetEX(ar.Context,key,value,expr).Err()
	if err != nil {
		ar.Logger.Error(fmt.Errorf("AuthRepository.SetRedis ERROR : %v MSG : %s",err,err.Error()))
		return err
	}
	return nil
} 

// GetRedis repository layer for get a value from redis by unique key
func (ar *AuthRepository) GetRedis(key string) (string,error) {
	cmd := ar.Redis.Get(ar.Context,key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return "", model.ErrRedisKeyNotExisted
		}

		ar.Logger.Error(fmt.Errorf("AuthRepository.GetRedis ERROR : %v MSG : %s",cmd.Err(),cmd.Err().Error()))

		return "",cmd.Err()
	}

	return cmd.String(),nil
}