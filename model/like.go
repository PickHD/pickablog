package model

type (

	// LikeRequest consist data for requesting create like
	LikeRequest struct {
		Like int `json:"like"`
		UserID int `json:"user_id"`
	}

	// ViewLikeResponse consist data of like response
	ViewLikeResponse struct {
		ID int `db:"id" json:"id"`
		Like int `db:"like_count" json:"like"`
		UserID int `db:"user_id" json:"user_id"`
		BlogID int `db:"article_id" json:"blog_id"`
	}
)