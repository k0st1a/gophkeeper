package client

import (
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
)

type Item struct {
	ID         int64
	Type       string
	Body       model.Item
	CreateTime time.Time
	UpdateTime time.Time
}
