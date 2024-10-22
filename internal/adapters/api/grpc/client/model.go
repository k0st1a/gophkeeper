package client

import (
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
)

// Item - модель Item клиента, для взаимодействия по GRPC с сервером.
type Item struct {
	// Тело предмета
	Body model.Item
	// Время создания предмета
	CreateTime time.Time
	// Время обновления предмета
	UpdateTime time.Time
	// Идентификатор предмета
	ID int64
}
