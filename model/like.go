package model

type (

	// LikeRequest consist data for requesting create like
	LikeRequest struct {
		Like int `json:"like"`
		UserID int `json:"user_id"`
	}
)