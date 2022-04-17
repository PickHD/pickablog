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
	// ICommentRepository is an interface that has all the function to be implemented inside comment repository
	ICommentRepository interface {
		Create(blogID int,req model.CommentRequest,createdBy string) error
		UpdateByID(id int,req map[string]interface{}, updatedBy string) error
		GetAllByBlogID(blogID int,page int, size int, order string, field string) ([]model.ViewCommentResponse,int,error)
		DeleteByID(blogID,commentID int) error
	}

	// CommentRepository is an app comment struct that consists of all the dependencies needed for comment repository
	CommentRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// Create repository layer for executing command to creating comments
func (cr *CommentRepository) Create(blogID int, req model.CommentRequest,createdBy string) error {
	tx,err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.Create BeginTX ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qInsert := `INSERT INTO comments (comment,article_id,user_id,created_by) VALUES ($1,$2,$3,$4) RETURNING id`

	var commentID int
	err = tx.QueryRow(cr.Context,qInsert,req.Comment,blogID,req.UserID,createdBy).Scan(&commentID)
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.Create.QueryRow.Scan Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		cr.Logger.Error(fmt.Errorf("CommentRepository.Create.QueryRow Scan ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qUpdate := `UPDATE article SET comments = ARRAY_APPEND(comments,$1) WHERE id = $2`

	_,err = tx.Exec(cr.Context,qUpdate,commentID,blogID)
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.Create.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		cr.Logger.Error(fmt.Errorf("CommentRepository.Create Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	err = tx.Commit(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.Create Commit ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// UpdateByID repository layer for executing command to updating a comment by id
func (cr *CommentRepository) UpdateByID(id int, req map[string]interface{}, updatedBy string) error {
	req["updated_by"] = updatedBy
	req["id"] = id

	tx,err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID BeginTX ERROR %v MSG %s",err,err.Error()))
		return err
	}

	q,args,err := helper.QueryUpdateBuilder("comments",req,[]string{"id"})
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID.QueryUpdateBuilder Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID QueryUpdateBuilder ERROR %v MSG %s",err,err.Error()))
		return err
	}

	cr.Logger.Info(fmt.Sprintf("Query : %s Args : %v",q,args))

	_,err = tx.Exec(cr.Context,q,args...)
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}


		cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	err = tx.Commit(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.UpdateByID Commit ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// GetAllByBlogID repository layer for querying command to getting all comments
func (cr *CommentRepository) GetAllByBlogID(blogID int,page int, size int, order string, field string) ([]model.ViewCommentResponse,int,error) { 
	q :=fmt.Sprintf(`
		SELECT
			id,
			comment,
			article_id,
			user_id,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM comments
		WHERE article_id = %d
	`,blogID)
	qCount := fmt.Sprintf(`SELECT 1 FROM comments WHERE article_id = %d`,blogID)

	limit := size + 1
	offset := (page - 1) * size
	orderBy := fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d ",field, order, limit, offset)

	query := fmt.Sprintf("%s %s", q, orderBy)
	queryCount := fmt.Sprintf("SELECT COUNT (*) FROM ( %s ) AS article_count ",qCount)

	cr.Logger.Info(fmt.Sprintf("Query : %s",query))
	cr.Logger.Info(fmt.Sprintf("Query Count : %s",queryCount))

	var totalData int
	err := cr.DB.QueryRow(cr.Context,queryCount).Scan(&totalData)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.GetAllByBlogID Scan ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	rows,err := cr.DB.Query(cr.Context,query)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.GetAllByBlogID Query ERROR %v MSG %s",err,err.Error()))
		return nil,0,err
	}

	var listData []model.ViewCommentResponse
	for rows.Next() {
		data := &model.ViewCommentResponse{}
		err := rows.Scan(&data.ID,&data.Comment,&data.BlogID,&data.UserID,&data.CreatedAt,&data.CreatedBy,&data.UpdatedAt,&data.UpdatedBy)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.GetAllByBlogID rows.Next Scan ERROR %v MSG %s",err,err.Error()))
			return nil,0,err
		}

		listData = append(listData, *data)
	}

	return listData,totalData,nil
}

// DeleteByID repository layer for executing command to deleting a comment by id
func (cr *CommentRepository) DeleteByID(blogID int,commentID int) error { 
	tx,err := cr.DB.Begin(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID BeginTX ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qDelete := `DELETE FROM comments WHERE id = $1`

	_,err = tx.Exec(cr.Context,qDelete,commentID)
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qUpdate := `UPDATE article SET comments = ARRAY_REMOVE(comments,$1) WHERE id = $2`

	_,err = tx.Exec(cr.Context,qUpdate,commentID,blogID)
	if err != nil {
		err = tx.Rollback(cr.Context)
		if err != nil {
			cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	err = tx.Commit(cr.Context)
	if err != nil {
		cr.Logger.Error(fmt.Errorf("CommentRepository.DeleteByID Commit ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}