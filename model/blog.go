package model

import (
	"database/sql"
	"regexp"
	"time"
)

var (
	// NoSpecialChar is regex validation where str is not contains a special characters
	NoSpecialChar *regexp.Regexp
)

type (
	// CreateBlogRequest consist data for creating a blog
	CreateBlogRequest struct {
		Title string `json:"title"`
		Slug string `json:"-"`
		Body string `json:"body"`
		Footer string `json:"footer"`
		UserID int `json:"user_id"`
		Tags []int `json:"tags"`
	}

	// UpdateBlogRequest consist data for update a blog
	UpdateBlogRequest struct {
		Title string `json:"title"`
		Slug string `json:"-"`
		Body string `json:"body"`
		Footer string `json:"footer"`
		UserID int `json:"-"`
	}

	// FilterBlogRequest consitss data for filter request blog
	FilterBlogRequest struct {
		StartDate string `query:"start_date"`
		EndDate string `query:"end_date"`
		Tags []int `query:"tags"`	
	}

	// ViewBlogResponse consists data of blog responses
	ViewBlogResponse struct {
		ID int `db:"id" json:"id"`
		Title string `db:"title" json:"title"`
		Body string `db:"body" json:"body"`
		Footer string `db:"footer" json:"footer"`
		UserID int `db:"user_id" json:"user_id"`
		Comments []sql.NullInt32 `db:"comments" json:"comments"`
		Tags []sql.NullInt32 `db:"tags" json:"tags"`
		Likes []sql.NullInt32 `db:"likes" json:"likes"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		CreatedBy string `db:"created_by" json:"created_by"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy *string `db:"updated_by" json:"updated_by,omitempty"`
	}

	BlogResponse struct {
		ID int `db:"id" json:"id"`
		Title string `db:"title" json:"title"`
		Body string `db:"body" json:"body"`
		Footer string `db:"footer" json:"footer"`
		UserID int `db:"user_id" json:"user_id"`
		Comments []int `db:"comments" json:"comments"`
		Tags []int `db:"tags" json:"tags"`
		Likes []int `db:"likes" json:"likes"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		CreatedBy string `db:"created_by" json:"created_by"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
		UpdatedBy *string `db:"updated_by" json:"updated_by,omitempty"`
	}
)

func init() {
	NoSpecialChar = regexp.MustCompile(`[^\w]`)
}