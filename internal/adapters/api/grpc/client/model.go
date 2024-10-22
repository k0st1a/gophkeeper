package client

import (
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
)

type Item struct {
	Body       model.Item
	CreateTime time.Time
	UpdateTime time.Time
	Type       string
	ID         int64
}
