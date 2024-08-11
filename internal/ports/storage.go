package ports

import (
	"context"
	"errors"
	"time"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login already busy")
	ErrUserNotFound     = errors.New("user not found")
)

type ItemStorage interface {
	CreateItem(ctx context.Context, userID int64, name, dataType string, data []byte) (int64, error)
	UpdateItem(ctx context.Context, userID int64, itemID int64, data []byte) error
	GetItem(ctx context.Context, userID int64, itemID int64) (*ItemInfo, error)
	ListItem(ctx context.Context, userID int64) ([]ItemInfo, error)
}

type ItemInfo struct {
	ID        int64
	Name      string
	Type      string
	Data      []byte
	CreatedAt time.Time
}

var ErrItemNotFound = errors.New("item not found")
