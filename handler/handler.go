package handler

import (
	"github.com/freekobie/kora/service"
)

type Handler struct {
	user *service.UserService
}

func NewHandler(us *service.UserService) *Handler {
	return &Handler{
		user: us,
	}
}
