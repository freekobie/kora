package handler

import (
	"github.com/freekobie/kora/model"
	"github.com/freekobie/kora/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func getUUIDparam(c *gin.Context, key string) (uuid.UUID, error) {
	idString := c.Param(key)
	return uuid.Parse(idString)
}
