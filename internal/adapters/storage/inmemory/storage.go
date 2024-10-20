package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/k0st1a/gophkeeper/internal/ports/client"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	mutex *sync.RWMutex
	items map[string]client.Item
}

func New() *Storage {
	return &Storage{
		mutex: &sync.RWMutex{},
		items: make(map[string]client.Item),
	}
}

func (s *Storage) Clear(ctx context.Context) {
	log.Printf("Clear")

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	s.items = make(map[string]client.Item)
}

// ListItems - возвращает копию списка предметов.
func (s *Storage) ListItems(ctx context.Context) ([]client.Item, error) {
	log.Printf("List items")

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return client.Map2List(s.items), nil
}

// GetItem - возвращает указатель на копию предмета.
func (s *Storage) GetItem(ctx context.Context, id string) (*client.Item, error) {
	log.Printf("Get item, id:%v", id)

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	i, ok := s.items[id]
	if !ok {
		log.Error().Msgf("Item(%v) not found", id)
		return nil, client.ErrorItemNotFound
	}

	return &i, nil
}

// CreateItem - создает предмет.
func (s *Storage) CreateItem(ctx context.Context, i *client.Item) (string, error) {
	log.Printf("Create item")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	i.ID = uuid.NewString()

	s.items[i.ID] = *i

	return i.ID, nil
}

// UpdateItem - обновляет предмет, если есть предмет с таким ID.
func (s *Storage) UpdateItem(ctx context.Context, ui *client.UpdateItem) error {
	log.Printf("Update item, id:%v", ui.ID)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	i, ok := s.items[ui.ID]
	if !ok {
		log.Error().Msgf("Item(%v) not found", ui.ID)
		return client.ErrorItemNotFound
	}

	updateItem(&i, ui)

	s.items[ui.ID] = i

	return nil
}

// UpdateItem - удаляет предмет.
func (s *Storage) DeleteItem(ctx context.Context, id string) error {
	log.Printf("Delete item, id:%v", id)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.items[id]
	if !ok {
		log.Error().Msgf("Item(%v) not found", id)
		return client.ErrorItemNotFound
	}

	delete(s.items, id)

	return nil
}

func updateItem(i *client.Item, ui *client.UpdateItem) {
	log.Printf("Start updateItem(%v)", i.ID)
	if ui.RemoteID != nil {
		log.Printf("Update(%v) RemoteID:%v", i.ID, *ui.RemoteID)
		i.RemoteID = *ui.RemoteID
	}

	if ui.Body != nil {
		log.Printf("Update(%v) Body", i.ID)
		i.Body = *ui.Body
	}

	if ui.UpdateTime != nil {
		log.Printf("Update(%v) UpdateTime:%v", i.ID, *ui.UpdateTime)
		i.UpdateTime = *ui.UpdateTime
	}

	if ui.DeleteMark != nil {
		log.Printf("Update(%v) DeleteMark:%v", i.ID, *ui.DeleteMark)
		i.DeleteMark = *ui.DeleteMark
	}
	log.Printf("End updateItem")
}
