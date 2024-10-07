package inmemory

import (
	"context"
	"sync"
	"time"

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

	return s.items
}

func (s *Storage) GetItem(ctx context.Context, Name string) (rawitem.Info, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	i, ok := s.items[Name]
	if !ok {
		return nil, rawitem.ErrorItemNotFound
	}

	return i, nil
}

func (s *Storage) AddItem(ctx context.Context, info rawitem.Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i, ok := s.items[info.Name]
	if ok {
		return rawitem.ErrorItemAlreadyExists
	}

	s.items[info.Name] = r

	return nil
}

func (s *Storage) UpdateItem(ctx context.Context, info rawitem.Info) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i, ok := s.items[Name]
	if !ok {
		return rawitem.ErrorItemNotFound
	}

	return nil
}

func (s *Storage) DeleteItem(ctx context.Context, Name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i, ok := s.items[Name]
	if !ok {
		return rawitem.ErrorItemNotFound
	}

	delete(s.items, Name)

	return nil
}
