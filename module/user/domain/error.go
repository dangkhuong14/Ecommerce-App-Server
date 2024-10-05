package domain

import "errors"

var (
	ErrEmailExists = errors.New("email already exists")
	ErrInvalidEmailPassword = errors.New("invalid email and password")
	ErrRefreshTokenExpired = errors.New("refresh token is expired")
	ErrUserBanned = errors.New("user is banned")
)