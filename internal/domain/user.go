package domain

import (
	"context"
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreateParameters struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(context.Context, UserCreateParameters) (User, error)
	GetByUsername(context.Context, string) (User, error)
	GetByID(context.Context, uint64) (User, error)
}
