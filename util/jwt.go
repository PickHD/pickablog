package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/model"
	"github.com/golang-jwt/jwt/v4"
)

var (
	// ErrSignedToken occurs when failed sign a jwt token
	ErrSignedToken = errors.New("failed sign a token %v")
	// ErrInvalidToken occurs when jwt token is invalid
	ErrInvalidToken = errors.New("invalid jwt token")
)


const (
	payloadFullName = "full_name"
	payloadEmail = "email"
	payloadExpires = "exp"
)

// BuildJWT return signed claims token JWT with defined expiration times in configuration 
func BuildJWT(cfg *config.Configuration, user *model.AuthUserDetails) (string,error) {
	claims := jwt.MapClaims{}
	claims[payloadFullName] = user.FullName
	claims[payloadEmail] = user.Email
	claims[payloadExpires] = time.Now().Add(time.Duration(cfg.Const.JWTExpire)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodES256,claims)

	signedToken,err := token.SignedString([]byte(cfg.Secret.JWTSecret))
	if err != nil {
		return "",fmt.Errorf(ErrSignedToken.Error(),err)
	}

	return signedToken,nil
}

// ValidateJWT will checking validity of signed JWT token
func ValidateJWT(cfg *config.Configuration,signedToken string) (*jwt.Token, error) {
	return jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.Secret.JWTSecret), nil
	})
}