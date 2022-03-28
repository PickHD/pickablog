package model

type (
	// AuthDetails consist data authorized users
	AuthUserDetails struct {
		FullName string `db:"full_name" json:"full_name"`
		Email string `db:"email" json:"email"`
		RoleID int `db:"id" json:"role_id"`
		RoleName string `db:"name" json:"role_name"`
	}

)

var (
	// KeyJWTValidAccess is context key identifier for valid jwt token
	KeyJWTValidAccess = "ValidJWTAccess"
)