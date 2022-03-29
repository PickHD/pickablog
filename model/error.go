package model

import "errors"

var (
	// ErrTokenNotFound occurs when jwt token not found
	ErrTokenNotFound = errors.New("token not found")
	// ErrSignedToken occurs when failed sign a jwt token
	ErrSignedToken = errors.New("failed sign a token %v")
	// ErrInvalidToken occurs when jwt token is invalid
	ErrInvalidToken = errors.New("invalid jwt token")
	// ErrTokenExpire occurs when jwt token already expired
	ErrTokenExpire = errors.New("token already expired, please to relogin application")
	
	// ErrTypeAssertion occurs when doing invalid type assertion
	ErrTypeAssertion = errors.New("type assertion error")

	// ErrFailedParseBody occurs when failed parsing body
	ErrFailedParseBody = errors.New("parse body error")

	// ErrInvalidRequets occurs when client sent invalid request body
	ErrInvalidRequest = errors.New("invalid request body")

	// ErrEmailExisted occurs when email already used / inside database
	ErrEmailExisted = errors.New("email is existed")
)