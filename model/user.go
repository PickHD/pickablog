package model

import "time"

type (
	// UpdateUserRequest consist data of updating a user
	UpdateUserRequest struct {
		FullName string `json:"full_name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// ViewUserResponse consist data of user
	ViewUserResponse struct {
		ID int `db:"id" json:"id,omitempty"`
		FullName string `db:"full_name" json:"full_name,omitempty"`
		Email string `db:"email" json:"email,omitempty"`
		CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
		CreatedBy string `db:"created_by" json:"created_by,omitempty"`
		UpdatedBy *string `db:"updated_by" json:"updated_by,omitempty"`
	}
)