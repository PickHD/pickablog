package model

import "time"

type (
	// CreateTagRequest consist data of creating a tag
	CreateTagRequest struct {
		Name string `json:"name"`
	}

	// UpdateTagRequest consist data of updating a tag
	UpdateTagRequest struct {
		Name string `json:"name"`
	}

	// ViewTagResponse consist data of tag
	ViewTagResponse struct {
		ID int `db:"id" json:"id,omitempty"`
		Name string `db:"name" json:"name,omitempty"`
		CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
		CreatedBy string `db:"created_by" json:"created_by,omitempty"`
		UpdatedBy *string `db:"updated_by" json:"updated_by,omitempty"`
	}
)