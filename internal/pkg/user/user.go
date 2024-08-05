package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/k0st1a/gophkeeper/internal/ports"
)

type Managment interface {
	Create(ctx context.Context, login, password string) (int64, error)
	GetIDAndPassword(ctx context.Context, login string) (int64, string, error)
}

type user struct {
	storage ports.UserStorage
}

var (
	ErrEmailAlreadyBusy = errors.New("user email already busy")
	ErrNotFound         = errors.New("user not found")
)

func New(storage ports.UserStorage) Managment {
	return &user{
		storage: storage,
	}
}

func (u *user) Create(ctx context.Context, email, password string) (int64, error) {
	id, err := u.storage.CreateUser(ctx, email, password)
	if err != nil {
		if errors.Is(err, ports.ErrEmailAlreadyBusy) {
			return 0, ErrEmailAlreadyBusy
		}

		return 0, fmt.Errorf("storage error of create user:%w", err)
	}

	return id, nil
}

func (u *user) GetIDAndPassword(ctx context.Context, email string) (int64, string, error) {
	id, password, err := u.storage.GetUserIDAndPassword(ctx, email)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return id, password, ErrNotFound
		}

		return id, password, fmt.Errorf("storage error of get user id and password:%w", err)
	}

	return id, password, nil
}
