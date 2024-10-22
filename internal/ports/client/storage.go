package client

import (
	"context"
	"errors"
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
)

type ItemStorage interface {
	Clear(ctx context.Context)
	CreateItem(ctx context.Context, item *Item) (string, error)
	UpdateItem(ctx context.Context, item *UpdateItem) error
	GetItem(ctx context.Context, id string) (*Item, error)
	ListItems(ctx context.Context) ([]Item, error)
	DeleteItem(ctx context.Context, id string) error
}

// Item - предмет, хранящегося в базе на стороне клиента.
type Item struct {
	// Тело предмета
	Body model.Item
	// Время создания предмета
	CreateTime time.Time
	// Время обновления предмета
	UpdateTime time.Time
	// Идентификатор предмета
	ID string
	// Идендификатор предмета на удаленном сервере
	RemoteID int64
	// Отметка о необходимости удаления предмета
	DeleteMark bool
}

// UpdateItem - для обновления полей предмета, хранящегося в базе на стороне клиента.
// Если поле не nil, значит его нужно обновить для предмета в хранилище.
type UpdateItem struct {
	// Идендификатор предмета на удаленном сервере
	RemoteID *int64
	// Тело предмета
	Body *model.Item
	// Идендификатор предмета на удаленном сервере
	UpdateTime *time.Time
	// Отметка о необходимости удаления предмета
	DeleteMark *bool
	// Идентификатор предмета
	ID string
}

var ErrItemNotFound = errors.New("item not found")

// List2MapWithRemoteID - преобразование из list в map, где ключом выступает RemoteID.
func List2MapWithRemoteID(l []*Item) map[int64]*Item {
	m := make(map[int64]*Item)
	for _, v := range l {
		m[v.RemoteID] = v
	}
	return m
}

// Map2List - преобразование из map в list.
func Map2List(m map[string]Item) []Item {
	l := make([]Item, len(m))
	i := 0
	for _, v := range m {
		l[i] = v
		i++
	}
	return l
}
