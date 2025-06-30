package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role,omitempty"`
	PasswordHash []byte    `json:"-"`
	ProfilePhoto string    `json:"profilePhoto"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModifed  time.Time `json:"lastModified"`
	Verified     bool      `json:"verified"`
}

type UserToken struct {
	Hash      string
	UserId    uuid.UUID
	ExpiresAt time.Time
	Scope     string
}

type UserStore interface {
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByMail(ctx context.Context, email string) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	InsertToken(ctx context.Context, token *UserToken) error
	GetUserForToken(ctx context.Context, tokenHash, scope, email string) (*User, error)
	DeleteToken(ctx context.Context, tokenHash, scope string) error
}
