package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type (
	// IBlogRepository is an interface that has all the function to be implemented inside blog repository
	IBlogRepository interface {
		Create(req model.CreateBlogRequest, createdBy string) error
		GetAll(page int, size int, order string, field string, search string,filter model.FilterBlogRequest) ([]model.ViewBlogResponse,int,error)
		GetBySlug(slug string) (*model.ViewBlogResponse,error)
		GetByID(id int) (*model.ViewBlogResponse,error)
		UpdateByID(id int, req map[string]interface{}, updatedBy string) error
		DeleteByID(id int) error
	}

	// BlogRepository is an app blog struct that consists of all the dependencies needed for blog repository
	BlogRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// Create repository layer for executing command create blog
func (br *BlogRepository) Create(req model.CreateBlogRequest, createdBy string) error {
	q := `INSERT INTO article (title,slug,body,footer,user_id,tags,created_by) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_,err := br.DB.Exec(br.Context,q,req.Title,req.Slug,req.Body,req.Footer,req.UserID,pq.Array(req.Tags),&createdBy)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.Create Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// GetAll repository layer for querying command get all blog
func (br *BlogRepository) GetAll(page int, size int, order string, field string, search string,filter model.FilterBlogRequest) ([]model.ViewBlogResponse,int,error) {
	q := `
	  	SELECT
		   	id,
			title,
			body,
			footer,
			user_id,
			comments,
			tags,
			likes,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM article
	`

	qCount := `SELECT 1 FROM article`

	criteria := ""
	criteria = ""

	if len(filter.StartDate) > 0 && len(filter.EndDate) > 0 {
		pStartDate,err := time.Parse(helper.StandardLayout,filter.StartDate)
		if err != nil {
			br.Logger.Error(fmt.Errorf("BlogRepository.GetAll time.Parse ERROR %v MSG %s",err,err.Error()))

			return nil,0,err
		}

		pEndDate,err := time.Parse(helper.StandardLayout,filter.EndDate)
		if err != nil {
			br.Logger.Error(fmt.Errorf("BlogRepository.GetAll time.Parse ERROR %v MSG %s",err,err.Error()))

			return nil,0,err
		}

		criteria += fmt.Sprintf(" created_at::DATE >= '%s' AND created_at::DATE <= '%s'",pStartDate.Format(helper.StandardLayout),pEndDate.Format(helper.StandardLayout))
	}

	tags := ""
	if len(filter.Tags) > 0 {
		tags = `tags @> '{`
		for i , t:= range filter.Tags {
			tags += fmt.Sprintf("%d",t)

			if i < len(filter.Tags) - 1 {
				tags += ","
			}
		}

		tags += "}'"

		if len(tags) > 0 {
			if len(criteria) > 0 {
				criteria += " AND " + tags
			} else {
				criteria += tags
			}
		}
	}

	if len(search) > 0 {
		if len(criteria) > 0 {
			criteria += " AND title ILIKE '%" + search + "%'"
		} else {
			criteria += " title ILIKE '%" + search + "%'"
		}
	} 

	cr := ""
	if len(criteria) > 0 {
		cr = " WHERE " + criteria
	}

	limit := size + 1
	offset := (page - 1) * size
	orderBy := fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d ",field, order, limit, offset)

	query := fmt.Sprintf("%s %s %s", q, cr, orderBy)
	queryCount := fmt.Sprintf("SELECT COUNT (*) FROM ( %s %s ) AS article_count ",qCount,cr)

	br.Logger.Info(fmt.Sprintf("Query : %s",query))
	br.Logger.Info(fmt.Sprintf("Query Count : %s",queryCount))

	var totalData int
	err := br.DB.QueryRow(br.Context,queryCount).Scan(&totalData)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.GetAll Scan ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	rows,err := br.DB.Query(br.Context,query)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.GetAll Query ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	var listData []model.ViewBlogResponse
	for rows.Next() {
		data := &model.ViewBlogResponse{}
		err := rows.Scan(
			&data.ID,
			&data.Title,
			&data.Body,
			&data.Footer,
			&data.UserID,
			pq.Array(&data.Comments),
			pq.Array(&data.Tags),
			pq.Array(&data.Likes),
			&data.CreatedAt,
			&data.CreatedBy,
			&data.UpdatedAt,
			&data.UpdatedBy,
		)
		if err != nil {
			br.Logger.Error(fmt.Errorf("BlogRepository.GetAll rows.Next Scan ERROR %v MSG %s",err,err.Error()))
			return nil,0,err
		} 

		listData = append(listData, *data)
	}

	
	return listData,totalData,nil
}

// GetBySlug repository layer for querying command get detail blog by Slug
func (br *BlogRepository) GetBySlug(slug string) (*model.ViewBlogResponse,error) {
	var blog model.ViewBlogResponse

	q := `
		SELECT 
			id,
			title,
			body,
			footer,
			user_id,
			comments,
			tags,
			likes,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM article
		WHERE slug = $1
	`

	row := br.DB.QueryRow(br.Context,q,slug)
	err := row.Scan(
		&blog.ID,
		&blog.Title,
		&blog.Body,
		&blog.Footer,
		&blog.UserID,
		pq.Array(&blog.Comments),
		pq.Array(&blog.Tags),
		pq.Array(&blog.Likes),
		&blog.CreatedAt,
		&blog.CreatedBy,
		&blog.UpdatedAt,
		&blog.UpdatedBy,
	)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.GetBySlug QueryRow.Scan ERROR %v MSG %s",err,err.Error()))
		return nil,err
	}

	return &blog,nil
}

// GetByID repository layer for querying command get detail blog by ID
func (br *BlogRepository) GetByID(id int) (*model.ViewBlogResponse,error) {
	var blog model.ViewBlogResponse

	q := `
		SELECT 
			id,
			title,
			body,
			footer,
			user_id,
			comments,
			tags,
			likes,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM article
		WHERE id = $1
	`

	row := br.DB.QueryRow(br.Context,q,id)
	err := row.Scan(
		&blog.ID,
		&blog.Title,
		&blog.Body,
		&blog.Footer,
		&blog.UserID,
		pq.Array(&blog.Comments),
		pq.Array(&blog.Tags),
		pq.Array(&blog.Likes),
		&blog.CreatedAt,
		&blog.CreatedBy,
		&blog.UpdatedAt,
		&blog.UpdatedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			br.Logger.Info(fmt.Errorf("BlogRepository.GetByID QueryRow.Scan INFO %v MSG %s",err,err.Error()))
		} else{
			br.Logger.Error(fmt.Errorf("BlogRepository.GetByID QueryRow.Scan ERROR %v MSG %s",err,err.Error()))
		}
		return nil,err
	}

	return &blog,nil
}

// UpdateByID repository layer for executing command update blog by ID
func (br *BlogRepository) UpdateByID(id int, req map[string]interface{},updatedBy string) error {
	req["updated_by"] = updatedBy
	req["id"] = id

	q,args,err := helper.QueryUpdateBuilder("article",req,[]string{"id"})
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.UpdateByID QueryUpdateBuilder ERROR %v MSG %s",err,err.Error()))
		return err
	}

	br.Logger.Info(fmt.Sprintf("Query : %s, Args : %v",q,args))

	_,err = br.DB.Exec(br.Context,q,args...)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.UpdateByID Exec ERROR %v MSG %s",err,err.Error()))

		return err
	}

	return nil
}

// DeleteByID repository layer for executin command delete blog by ID
func (br *BlogRepository) DeleteByID(id int) error {
	q := `DELETE FROM article WHERE id = $1`

	_,err := br.DB.Exec(br.Context,q,id)
	if err != nil {
		br.Logger.Error(fmt.Errorf("BlogRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))

		return err
	}

	return nil
}