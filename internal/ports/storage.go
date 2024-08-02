package ports

import (
	"context"
	"errors"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login already busy")
	ErrUserNotFound     = errors.New("user not found")
)
