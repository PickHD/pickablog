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
	// ITagRepository is an interface that has all the function to be implemented inside tag repository
	ITagRepository interface {
		Create(req model.CreateTagRequest, createdBy string) error
		GetByName(name string) (*model.ViewTagResponse,error)
		GetAll(page int, size int, order string, field string, search string) ([]model.ViewTagResponse,int,error)
		GetByID(id int) (*model.ViewTagResponse,error)
		UpdateByID(id int,req model.UpdateTagRequest,updatedBy string) error
		DeleteByID(id int) error
	}

	// TagRepository is an app health check struct that consists of all the dependencies needed for tag repository
	TagRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// Create repository layer for executing command create a tag
func (tr *TagRepository) Create(req model.CreateTagRequest,createdBy string) error {
	q := `INSERT INTO tag (name,created_by) VALUES ($1,$2)`

	_, err := tr.DB.Exec(tr.Context,q,req.Name,createdBy)
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.Create Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// GetByName repository layer for quering command getting tag by name
func (tr *TagRepository) GetByName(name string) (*model.ViewTagResponse,error) {
	var tag model.ViewTagResponse

	q := ` 
		SELECT 
			id,
			name,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM tag
		WHERE name = $1
	`

	row := tr.DB.QueryRow(tr.Context,q,name)
	err := row.Scan(&tag.ID,&tag.Name,&tag.CreatedAt,&tag.UpdatedAt,&tag.CreatedBy,&tag.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			tr.Logger.Info(fmt.Errorf("TagRepository.GetByName Scan INFO %v MSG %s",err,err.Error()))
		} else {
			tr.Logger.Error(fmt.Errorf("TagRepository.GetByName Scan ERROR %v MSG %s",err,err.Error()))
		}

		return nil,err
	}

	return &tag,nil
}

// GetAll repository layer for quering command getting all tag
func (tr *TagRepository) GetAll(page int,size int,order string, field string,search string) ([]model.ViewTagResponse,int,error) {
	q := `
		SELECT 
			id,
			name,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM tag
	`
	qCount := `SELECT 1 FROM tag`

	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " name LIKE '%" + search + "%'"
	}

	cr := ""
	if len(criteria) > 0 {
		cr = " WHERE " + criteria
	}

	limit := size + 1
	offset := (page - 1) * size
	orderBy := fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d ",field, order, limit, offset)

	query := fmt.Sprintf("%s %s %s", q, cr, orderBy)
	queryCount := fmt.Sprintf("SELECT COUNT (*) FROM ( %s %s ) AS tag_count ",qCount,cr)

	tr.Logger.Info(fmt.Sprintf("Query : %s",query))
	tr.Logger.Info(fmt.Sprintf("Query Count : %s",queryCount))

	var totalData int
	err := tr.DB.QueryRow(tr.Context,queryCount).Scan(&totalData)
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.GetAll Scan ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	rows,err := tr.DB.Query(tr.Context,query)
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.GetAll Query ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	var listData []model.ViewTagResponse
	for rows.Next() {
		data := &model.ViewTagResponse{}
		err := rows.Scan(&data.ID,&data.Name,&data.CreatedAt,&data.UpdatedAt,&data.CreatedBy,&data.UpdatedBy)
		if err != nil {
			tr.Logger.Error(fmt.Errorf("TagRepository.GetAll rows.Next Scan ERROR %v MSG %s",err,err.Error()))
			return nil,0,err
		}

		listData = append(listData,*data)
	}

	return listData,totalData,nil
}

// GetByID repository layer for quering command getting tag by id
func (tr *TagRepository) GetByID(id int) (*model.ViewTagResponse,error) {
	var tag model.ViewTagResponse

	q := ` 
		SELECT 
			id,
			name,
			created_at,
			updated_at,
			created_by,
			updated_by
		FROM tag
		WHERE id = $1
	`

	row := tr.DB.QueryRow(tr.Context,q,id)
	err := row.Scan(&tag.ID,&tag.Name,&tag.CreatedAt,&tag.UpdatedAt,&tag.CreatedBy,&tag.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			tr.Logger.Info(fmt.Errorf("TagRepository.GetByID Scan INFO %v MSG %s",err,err.Error()))
		} else {
			tr.Logger.Error(fmt.Errorf("TagRepository.GetByID Scan ERROR %v MSG %s",err,err.Error()))
		}

		return nil,err
	}

	return &tag,nil
}

// UpdateByID repository layer for executing command updating tag by id
func (tr *TagRepository) UpdateByID(id int, req model.UpdateTagRequest,updatedBy string) error {
	fields := make(map[string]interface{})
	fields["name"] = req.Name
	fields["updated_by"] = updatedBy
	fields["id"] = id

	q,args,err := helper.QueryUpdateBuilder("tag",fields,[]string{"id"})
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.UpdateByID QueryUpdateBuilder ERROR %v MSG %s",err,err.Error()))
		return err
	}

	tr.Logger.Info(fmt.Sprintf("Query : %s Args : %v",q,args))

	_,err = tr.DB.Exec(tr.Context,q,args...)
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.UpdateByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// DeleteByID repository layer for executing command deleting tag by id
func (tr *TagRepository) DeleteByID(id int) error {
	q := `DELETE FROM tag WHERE id = $1`

	_,err := tr.DB.Exec(tr.Context,q,id)
	if err != nil {
		tr.Logger.Error(fmt.Errorf("TagRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}