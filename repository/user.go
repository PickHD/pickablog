package repository

import (
	"context"
	"fmt"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IUserRepository is an interface that has all the function to be implemented inside user repository
	IUserRepository interface {
		GetAll(page int, size int, order string, field string, search string) ([]model.ViewUserResponse,int,error)
		GetByEmail(email string) (*model.ViewUserResponse,error)
		GetByID(id int) (*model.ViewUserResponse,error)
		UpdateByID(id int, req map[string]interface{}, updatedBy string) error
		DeleteByID(id int) error
	}

	// UserRepository is an app user check struct that consists of all the dependencies needed for user repository
	UserRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// GetAll repository layer for querying command getting all user
func (ur *UserRepository) GetAll(page int, size int, order string, field string, search string) ([]model.ViewUserResponse,int,error) {
	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM "user"
		WHERE role_id = 2
	`
	qCount := `SELECT 1 FROM "user" WHERE role_id = 2`

	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " AND full_name LIKE '%" + search + "%' OR email LIKE '%" + search + "%'"
	}

	limit := size + 1
	offset := (page - 1) * size
	orderBy := fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d ",field, order, limit, offset)

	query := fmt.Sprintf("%s %s %s", q, criteria, orderBy)
	queryCount := fmt.Sprintf("SELECT COUNT (*) FROM ( %s %s ) AS user_count ",qCount,criteria)

	ur.Logger.Info(fmt.Sprintf("Query : %s",query))
	ur.Logger.Info(fmt.Sprintf("Query Count : %s",queryCount))
 
 	var totalData int
	err := ur.DB.QueryRow(ur.Context,queryCount).Scan(&totalData)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.GetAll Scan ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	rows,err := ur.DB.Query(ur.Context,query)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.GetAll Query ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	var listData []model.ViewUserResponse
	for rows.Next() {
		data := &model.ViewUserResponse{}
		err := rows.Scan(&data.ID,&data.FullName,&data.Email,&data.CreatedAt,&data.UpdatedAt,&data.CreatedBy,&data.UpdatedBy)
		if err != nil {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetAll rows.Next Scan ERROR %v MSG %s",err,err.Error()))
			return nil,0,err
		}

		listData = append(listData, *data)
	}

	return listData,totalData,nil
}

// GetByEmail repository layer for querying command get a user by email
func (ur *UserRepository) GetByEmail(email string) (*model.ViewUserResponse,error) {
	var user model.ViewUserResponse
	
	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM "user"
		WHERE email = $1
	`

	row := ur.DB.QueryRow(ur.Context,q,email)
	err := row.Scan(&user.ID,&user.FullName,&user.Email,&user.CreatedAt,&user.UpdatedAt,&user.CreatedBy,&user.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.Logger.Info(fmt.Errorf("UserRepository.GetByEmail Scan INFO %v MSG %s",err,err.Error()))
		} else {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetByEmail Scan ERROR %v MSG %s",err,err.Error()))
		}

		return nil,err
	}

	return &user,nil
}

// GetByID repository layer for querying command get a user by id
func (ur *UserRepository) GetByID(id int) (*model.ViewUserResponse,error) {
	var user model.ViewUserResponse
	
	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM "user"
		WHERE id = $1
	`

	row := ur.DB.QueryRow(ur.Context,q,id)
	err := row.Scan(&user.ID,&user.FullName,&user.Email,&user.CreatedAt,&user.UpdatedAt,&user.CreatedBy,&user.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.Logger.Info(fmt.Errorf("UserRepository.GetByID Scan INFO %v MSG %s",err,err.Error()))
		} else {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetByID Scan ERROR %v MSG %s",err,err.Error()))
		}

		return nil,err
	}

	return &user,nil
}

// UpdateByID repository layer for executing command update a user by id
func (ur *UserRepository) UpdateByID(id int, req map[string]interface{},updatedBy string) error {
	req["updated_by"] = updatedBy
	req["id"] = id

	q,args,err := helper.QueryUpdateBuilder(`"user"`,req,[]string{"id"})
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.UpdateByID QueryUpdateBuilder ERROR %v MSG %s",err,err.Error()))
		return err
	}

	ur.Logger.Info(fmt.Sprintf("Query : %s Args : %v",q,args))

	_,err = ur.DB.Exec(ur.Context,q,args...)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.UpdateByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// UpdateByID repository layer for executing command delete a user by id
func (ur *UserRepository) DeleteByID(id int) error {
	q := `DELETE FROM "user" WHERE id = $1`

	_,err := ur.DB.Exec(ur.Context,q,id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}