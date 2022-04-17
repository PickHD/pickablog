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
	// ILikeRepository is an interface that has all the function to be implemented inside like repository
	ILikeRepository interface {
		Create(blogID int,req model.LikeRequest,createdBy string) error
		DeleteByID(blogID int,likeID int) error
	}

	// LikeRepository is an app like struct that consists of all the dependencies needed for like repository
	LikeRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// Create repository layer for executing command creating like 
func (lr *LikeRepository) Create(blogID int,req model.LikeRequest,createdBy string) error {
	tx,err := lr.DB.Begin(lr.Context)
	if err != nil {
		lr.Logger.Error(fmt.Errorf("LikeRepository.Create BeginTX ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qInsert := `INSERT INTO likes (like_count,article_id,user_id,created_by) VALUES ($1,$2,$3,$4) RETURNING id`

	var likeID int
	err = tx.QueryRow(lr.Context,qInsert,req.Like,blogID,req.UserID,createdBy).Scan(&likeID)
	if err != nil {
		err = tx.Rollback(lr.Context)
		if err != nil {
			lr.Logger.Error(fmt.Errorf("LikeRepository.Create.QueryRow.Scan Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		lr.Logger.Error(fmt.Errorf("LikeRepository.Create.QueryRow Scan ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qUpdate := `UPDATE article SET likes = ARRAY_APPEND(likes,$1) WHERE id = $2`

	_,err = tx.Exec(lr.Context,qUpdate,likeID,blogID)
	if err != nil {
		err = tx.Rollback(lr.Context)
		if err != nil {
			lr.Logger.Error(fmt.Errorf("LikeRepository.Create.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		lr.Logger.Error(fmt.Errorf("LikeRepository.Create Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	err = tx.Commit(lr.Context)
	if err != nil {
		lr.Logger.Error(fmt.Errorf("LikeRepository.Create Commit ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}

// DeleteByID repository layer for executing command deleting like by id
func (lr *LikeRepository) DeleteByID(blogID int,likeID int) error {
	tx,err := lr.DB.Begin(lr.Context)
	if err != nil {
		lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID BeginTX ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qDelete := `DELETE FROM likes WHERE id = $1`

	_,err = tx.Exec(lr.Context,qDelete,likeID)
	if err != nil {
		err = tx.Rollback(lr.Context)
		if err != nil {
			lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	qUpdate := `UPDATE article SET comments = ARRAY_REMOVE(comments,$1) WHERE id = $2`

	_,err = tx.Exec(lr.Context,qUpdate,likeID,blogID)
	if err != nil {
		err = tx.Rollback(lr.Context)
		if err != nil {
			lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID.Exec Rollback ERROR %v MSG %s",err,err.Error()))
			return err
		}

		lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID Exec ERROR %v MSG %s",err,err.Error()))
		return err
	}

	err = tx.Commit(lr.Context)
	if err != nil {
		lr.Logger.Error(fmt.Errorf("LikeRepository.DeleteByID Commit ERROR %v MSG %s",err,err.Error()))
		return err
	}

	return nil
}