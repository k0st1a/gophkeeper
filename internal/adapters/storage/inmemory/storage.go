package inmemory

import (
	"context"
	"sync"

	"github.com/k0st1a/gophkeeper/internal/pkg/rawitem"
)

type Storage struct {
	mutex *sync.RWMutex
	items map[string]rawitem.Info
}

func New() *Storage {
	return &Storage{
		mutex: &sync.RWMutex{},
		items: make(map[string]rawitem.Info),
	}
}

func (s *Storage) ListItems(ctx context.Context) []rawitem.Info {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return rawitem.Map2List(s.items)
}

// GetItem - возвращает указатель на копию предмета из Storage.
func (s *Storage) GetItem(ctx context.Context, Name string) (*rawitem.Info, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	i, ok := s.items[Name]
	if !ok {
		return nil, rawitem.ErrorItemNotFound
	}

	info := i

	return &info, nil
}

// AddItem - добавляет предмет в Storage.
func (s *Storage) AddItem(ctx context.Context, info *rawitem.Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.items[info.Name]
	if ok {
		return rawitem.ErrorItemAlreadyExists
	}

	s.items[info.Name] = *info

	return nil
}

// UpdateItem - обновляет предмет в Storage.
func (s *Storage) UpdateItem(ctx context.Context, info *rawitem.Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.items[info.Name]
	if !ok {
		return rawitem.ErrorItemNotFound
	}

	s.items[info.Name] = *info

	return nil
}

// UpdateItem - удаляет предмет из Storage.
func (s *Storage) DeleteItem(ctx context.Context, Name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.items[Name]
	if !ok {
		return rawitem.ErrorItemNotFound
	}

	delete(s.items, Name)

	return nil
}
