package handler

import (
	"github.com/freekobie/kora/model"
	"github.com/freekobie/kora/session"
)

type UserResponse struct {
	Status int        `json:"status"`
	User   model.User `json:"user"`
}

type SessionResponse struct {
	Status  int                 `json:"status"`
	Session session.UserSession `json:"session"`
}

type AccessResponse struct {
	Status int                `json:"status"`
	Access session.UserAccess `json:"access"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
