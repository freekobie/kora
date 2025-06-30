package handler

import "errors"

var (
	ErrServerError = errors.New("the server could not process your request")
)
