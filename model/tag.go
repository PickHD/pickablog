package model

import "time"

type (
	// CreateTagRequest consist data of creating a tag
	CreateTagRequest struct {
		Name string `json:"name"`
	}

	ViewTagResponse struct {
		ID int `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}
)