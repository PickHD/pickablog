package model

import "time"

type (

	// CommentRequest consist data for requesting create/update comment
	CommentRequest struct {
		Comment string `json:"comment"`
		UserID int `json:"user_id"`
	}

	// ViewCommentResponse consist data of comments
	ViewCommentResponse struct {
		ID int `db:"id" json:"id"`
		Comment string `db:"comment" json:"comment"`
		BlogID int `db:"article_id" json:"blog_id"`
		UserID int `db:"user_id" json:"user_id"`
		CreatedAt time.Time	`db:"created_at" json:"created_at"`
		CreatedBy string `db:"created_by" json:"created_by"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy *string `db:"updated_by" json:"updated_by,omitempty"`
	}
)