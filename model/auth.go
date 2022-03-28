package model

type (
	// AuthDetails consist data authorized users
	AuthUserDetails struct {
		FullName string `db:"full_name" json:"full_name"`
		Email string `db:"email" json:"email"`
		RoleID int `db:"role_id" json:"role_id"`
	}
)