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
	// ITagRepository is an interface that has all the function to be implemented inside tag repository
	ITagRepository interface {
		Create(req model.CreateTagRequest, createdBy string) error
		GetByName(name string) (*model.ViewTagResponse,error)
		GetAll(page int, size int, order string, field string, search string) ([]model.ViewTagResponse,int,error)
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
			updated_at 
		FROM tag
		WHERE name = $1
	`

	row := tr.DB.QueryRow(tr.Context,q,name)
	err := row.Scan(&tag.ID,&tag.Name,&tag.CreatedAt,&tag.UpdatedAt)
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
			updated_at 
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
		err := rows.Scan(&data.ID,&data.Name,&data.CreatedAt,&data.UpdatedAt)
		if err != nil {
			tr.Logger.Error(fmt.Errorf("TagRepository.GetAll rows.Next Scan ERROR %v MSG %s",err,err.Error()))
			return nil,0,err
		}

		listData = append(listData,*data)
	}

	return listData,totalData,nil
}