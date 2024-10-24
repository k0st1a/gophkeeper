package storage

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
	pclient "github.com/k0st1a/gophkeeper/internal/ports/client"
	"github.com/rs/zerolog/log"
)

type ItemStorage interface {
	Clear(ctx context.Context)
	CreateItem(ctx context.Context, body any) (string, error)
	UpdateItem(ctx context.Context, item *Item) error
	GetItem(ctx context.Context, id string) (*Item, error)
	ListItems(ctx context.Context) ([]Item, error)
	DeleteItem(ctx context.Context, id string) error
}

type client struct {
	storage pclient.ItemStorage
}

// New - создать клиента для работы с хранилищем предметов.
func New(s pclient.ItemStorage) *client {
	return &client{
		storage: s,
	}
}

// Clear – подчистка хранилища.
func (c *client) Clear(ctx context.Context) {
	c.storage.Clear(ctx)
}

// CreateItem - создать предмет.
func (c *client) CreateItem(ctx context.Context, body any) (string, error) {
	log.Ctx(ctx).Printf("Create item")

	mi, err := createStorageItemBody(body)
	if err != nil {
		return "", err
	}

	now := time.Now()

	si := &pclient.Item{
		Body:       *mi,
		CreateTime: now,
		UpdateTime: now,
		DeleteMark: false,
	}

	id, err := c.storage.CreateItem(ctx, si)
	if err != nil {
		return "", fmt.Errorf("error of create item:%w", err)
	}

	log.Ctx(ctx).Printf("Item created, id:%v", id)
	return id, nil
}

// UpdateItem - обновить предмет.
func (c *client) UpdateItem(ctx context.Context, i *Item) error {
	log.Ctx(ctx).Printf("Update item(%v)", i.ID)

	mi, err := createStorageItemBody(i.Body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of create storage item(%v) body", i.ID)
		return fmt.Errorf("error of create storage item(%v) body:%w", i.ID, err)
	}

	ut := time.Now()

	ui := &pclient.UpdateItem{
		ID:         i.ID,
		Body:       mi,
		UpdateTime: &ut,
	}

	err = c.storage.UpdateItem(ctx, ui)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of update storage item(%v)", i.ID)
		return fmt.Errorf("error of update storage item(%v):%w", i.ID, err)
	}

	log.Ctx(ctx).Printf("Item(%v) updated", i.ID)
	return nil
}

// GetItem - получить предмет по его идентификатору.
func (c *client) GetItem(ctx context.Context, id string) (*Item, error) {
	log.Ctx(ctx).Printf("Get item(%v)", id)

	si, err := c.storage.GetItem(ctx, id)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of get storage item(%v)", id)
		return nil, fmt.Errorf("error of get item(%v):%w", id, err)
	}

	if si.DeleteMark {
		return nil, fmt.Errorf("Item(%v) mark to delete", id)
	}

	b, err := createItemBody(&si.Body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of create item(%v) body", id)
		return nil, fmt.Errorf("error of create item(%v) body:%w", id, err)
	}

	i := &Item{
		ID:         si.ID,
		Body:       b,
		CreateTime: si.CreateTime,
		UpdateTime: si.UpdateTime,
	}

	log.Ctx(ctx).Printf("Item(%v) got", id)
	return i, nil
}

// ListItems - получить список предметов.
func (c *client) ListItems(ctx context.Context) ([]Item, error) {
	log.Ctx(ctx).Printf("List items")

	sl, err := c.storage.ListItems(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error of list items")
		return nil, fmt.Errorf("error of list items:%w", err)
	}

	l := make([]Item, 0, len(sl))

	for _, si := range sl {
		if si.DeleteMark {
			continue
		}

		b, err := createItemBody(&si.Body)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error of create item(%v) body", si.ID)
			continue
		}

		i := Item{
			ID:         si.ID,
			Body:       b,
			CreateTime: si.CreateTime,
			UpdateTime: si.UpdateTime,
		}

		l = append(l, i)
	}

	log.Ctx(ctx).Printf("List items success")
	return l, nil
}

// DeleteItem - удалить предмет по его идентификатору.
func (c *client) DeleteItem(ctx context.Context, id string) error {
	log.Ctx(ctx).Printf("Delete item(%v)", id)

	dm := true
	ut := time.Now()

	ui := pclient.UpdateItem{
		ID:         id,
		DeleteMark: &dm,
		UpdateTime: &ut,
	}

	err := c.storage.UpdateItem(ctx, &ui)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of update item(%v) while mark to delete", id)
		return fmt.Errorf("error of get item(%v):%w", id, err)
	}

	log.Ctx(ctx).Printf("Item(%v) mark to delete", id)
	return nil
}

func createStorageItemBody(body any) (*model.Item, error) {
	var mi model.Item

	switch b := body.(type) {
	case *Password:
		mi.Password = &model.Password{
			Resource: b.Resource,
			UserName: b.UserName,
			Password: b.Password,
		}
	case *Card:
		mi.Card = &model.Card{
			Number:  b.Number,
			Expires: b.Expires,
			Holder:  b.Holder,
		}
	case *Note:
		mi.Note = &model.Note{
			Name: b.Name,
			Body: b.Body,
		}
	case *File:
		mi.File = &model.File{
			Name:        b.Name,
			Description: b.Description,
			Body:        b.Body,
		}
	default:
		return nil, fmt.Errorf("unkown item body type:%v", reflect.TypeOf(b))
	}

	return &mi, nil
}

func createItemBody(mi *model.Item) (any, error) {
	mib, err := mi.GetBody()
	if err != nil {
		return nil, fmt.Errorf("error of get item body:%w", err)
	}

	var ib any

	switch b := mib.(type) {
	case *model.Password:
		ib = &Password{
			Resource: b.Resource,
			UserName: b.UserName,
			Password: b.Password,
		}
	case *model.Card:
		ib = &Card{
			Number:  b.Number,
			Expires: b.Expires,
			Holder:  b.Holder,
		}
	case *model.Note:
		ib = &Note{
			Name: b.Name,
			Body: b.Body,
		}
	case *model.File:
		ib = &File{
			Name:        b.Name,
			Description: b.Description,
			Body:        b.Body,
		}
	default:
		return nil, fmt.Errorf("unkown storage item body:%v", reflect.TypeOf(b))
	}

	return ib, nil
}
