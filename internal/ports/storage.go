package ports

import (
	"context"
	"errors"
)

type UserStorage interface {
	CreateUser(ctx context.Context, email, password string) (int64, error)
	GetUserIDAndPassword(ctx context.Context, email string) (int64, string, error)
}

var (
	ErrEmailAlreadyBusy = errors.New("email already busy")
	ErrUserNotFound     = errors.New("user not found")
)
