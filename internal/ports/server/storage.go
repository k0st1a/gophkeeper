package server

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
	CreateItem(ctx context.Context, userID int64, item *Item) (int64, error)
	UpdateItem(ctx context.Context, userID int64, item *Item) error
	GetItem(ctx context.Context, userID int64, itemID int64) (*Item, error)
	ListItems(ctx context.Context, userID int64) ([]Item, error)
	DeleteItem(ctx context.Context, userID int64, itemID int64) error
}

type Item struct {
	CreateTime time.Time
	UpdateTime time.Time
	Data       []byte
	ID         int64
}

var ErrItemNotFound = errors.New("item not found")
