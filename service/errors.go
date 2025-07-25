package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnverifiedUser     = errors.New("user has an unverified email")
	ErrInvalidToken       = errors.New("token is invalid or expired")
	ErrFailedOperation    = errors.New("failed to complete operation")
	ErrInvalidPassword    = errors.New("password must be between 8 and 20 characters")
)
