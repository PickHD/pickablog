package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
	"github.com/PickHD/pickablog/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type (
	DecodePayloadData struct {
		FullName string `json:"full_name"`
		Email string `json:"email"`
		RoleName string `json:"role_name"`
	}
)

const (
	payloadFullName string = "full_name"
	payloadEmail string = "email"
	payloadRoleName string = "role_name"
	payloadExpires string = "exp"

	JWTExpire time.Duration = time.Duration(7) * time.Hour
)

// BuildJWT return signed claims token JWT with defined expiration times in configuration 
func BuildJWT(cfg *config.Configuration, user *model.AuthUserDetails) (string,error) {
	claims := jwt.MapClaims{}
	claims[payloadFullName] = user.FullName
	claims[payloadEmail] = user.Email
	claims[payloadRoleName] = user.RoleName
	claims[payloadExpires] = time.Now().Add(JWTExpire).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	signedToken,err := token.SignedString([]byte(cfg.Secret.JWTSecret))
	if err != nil {
		return "",fmt.Errorf(model.ErrSignedToken.Error(),err)
	}

	return signedToken,nil
}

// ValidateJWT will checking validity of signed JWT token from request in
func ValidateJWT(ctx *fiber.Ctx) (DecodePayloadData, error) {
	header := ctx.Get("Authorization","")
	if !strings.Contains(header,"Bearer") {
		return DecodePayloadData{},model.ErrTokenNotFound
	}

	getToken := strings.Replace(header, "Bearer ", "", -1)
	validToken, err := jwt.Parse(getToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, model.ErrInvalidToken
		}
		return []byte(helper.GetEnvString("JWT_SECRET")), nil
	})
	if err != nil {
		return DecodePayloadData{},model.ErrInvalidToken
	}

	claims := validToken.Claims.(jwt.MapClaims)

	// Check is token expired or not
	if expInt, ok := claims[payloadExpires].(float64); ok {
		now := time.Now().Unix()
		if now > int64(expInt) {
			return DecodePayloadData{},model.ErrTokenExpire
		}
	} else {
		return DecodePayloadData{},model.ErrTypeAssertion
	}

	decodePayload,err := insertPayloadJWT(claims)
	if err != nil {
		return DecodePayloadData{}, err
	}

	return decodePayload,nil
}

// ExtractPayloadJWT will extracting payload data from ctx.Locals
func ExtractPayloadJWT(data interface{}) (DecodePayloadData, error) {
	extractData,ok := data.(DecodePayloadData)
	if !ok {
		return DecodePayloadData{},model.ErrTypeAssertion
	}

	return extractData,nil
}

// insertPayloadJWT will inserting data from decoded payload into defined struct
func insertPayloadJWT(claims jwt.MapClaims) (DecodePayloadData,error) {
	decodePayloadData := DecodePayloadData{}

	if userEmail, ok := claims[payloadEmail].(string); ok {
		decodePayloadData.Email = userEmail
	} else {
		return DecodePayloadData{}, model.ErrTypeAssertion
	}

	if userFullName, ok := claims[payloadFullName].(string); ok {
		decodePayloadData.FullName = userFullName
	} else {
		return DecodePayloadData{}, model.ErrTypeAssertion
	}

	if userRoleName,ok := claims[payloadRoleName].(string); ok {
		decodePayloadData.RoleName = userRoleName
	} else {
		return DecodePayloadData{}, model.ErrTypeAssertion
	}

	return decodePayloadData,nil
}