package model

import (
	"regexp"
	"time"
)

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
		Password string `db:"password" json:"-"`
	}

	// CreateUserRequest consist data for creating a user
	CreateUserRequest struct {
		FullName string `json:"full_name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// LoginRequest consist data for log-in a user
	LoginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// GoogleOauthResponse consist data from response google services user info 
	GoogleOauthResponse struct {
		ID string `json:"id"`
		Email string `json:"email"`
		VerifiedEmail bool `json:"verified_email"`
		Name string `json:"name,omitempty"`
		GivenName string `json:"given_name,omitempty"`
		Picture string `json:"picture,omitempty"`
		Locale string `json:"locale"`
		Time string `json:"time"`
	}
	// SuccessLoginResponse consist data of success login 
	SuccessLoginResponse struct {
		AccessToken string `json:"access_token"`
		ExpiredAt time.Time `json:"expired_at"`
		Role string `json:"role"`
	}

)

const (
	// RoleAuthor mapping id of role author
	RoleAuthor int = 2
	
	// RoleGuest mapping id of role guest
	RoleGuest int = 3

	//OauthStateKey consists template key of oauth state for redis
	OauthStateKey string = "oauth_state"
)

func init() {
	IsAllowedEmailInput = regexp.MustCompile(`^[^\s@]+@([^\s@.,]+\.)+[^\s@.,]{2,}$`)
}

// GetValidRoleByID will return valid role by its id
func GetValidRoleByID(roleID int) (string,error) {
	switch roleID {
	case RoleAuthor:
		return "Author",nil
	case RoleGuest:
		return "Guest",nil
	default:
		return "",ErrRoleNotExisted
	}
}