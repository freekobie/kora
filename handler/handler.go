package handler

import (
	"github.com/freekobie/kora/service"
)

type Handler struct {
	user *service.UserService
	file *service.FileService
}

func NewHandler(us *service.UserService, fs *service.FileService) *Handler {
	return &Handler{
		user: us,
		file: fs,
	}
}
