package model

import "regexp"

var (
	// KeyJWTValidAccess is context key identifier for valid jwt token
	KeyJWTValidAccess = "ValidJWTAccess"

	// IsAllowedEmailInput is regex validator to allowing only valid email
	IsAllowedEmailInput *regexp.Regexp
)

type (
	// AuthDetails consist data authorized users
	AuthUserDetails struct {
		FullName string `db:"full_name" json:"full_name"`
		Email string `db:"email" json:"email"`
		RoleID int `db:"id" json:"role_id"`
		RoleName string `db:"name" json:"role_name"`
	}

	// RegisterAuthorRequest consist data for register as a author
	RegisterAuthorRequest struct {
		FullName string `json:"full_name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

)

func init() {
	IsAllowedEmailInput = regexp.MustCompile(`^[^\s@]+@([^\s@.,]+\.)+[^\s@.,]{2,}$`)
}