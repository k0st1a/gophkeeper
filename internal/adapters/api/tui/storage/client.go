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
	CreateItem(ctx context.Context, body any, meta Meta) (string, error)
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
func (c *client) CreateItem(ctx context.Context, body any, meta Meta) (string, error) {
	log.Ctx(ctx).Printf("Create item")

	var item model.Item

	err := convertAndFillBody(&item, body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of convert and fill body while create item")
		return "", fmt.Errorf("error of convert and fill body while create item:%w", err)
	}
	item.Meta = model.Meta(meta)

	now := time.Now()

	si := &pclient.Item{
		Body:       item,
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

	var item model.Item

	err := convertAndFillBody(&item, i.Body)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("error of convert and fill body while update item(%v)", i.ID)
		return fmt.Errorf("error of convert and fill body while update item(%v):%w", i.ID, err)
	}
	item.Meta = model.Meta(i.Meta)

	ut := time.Now()
	ui := &pclient.UpdateItem{
		ID:         i.ID,
		Body:       &item,
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

	i, err := parseItem(si)
	if err != nil {
		return nil, fmt.Errorf("error of parse item(%v) while get item", id)
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

		i, err := parseItem(&si)
		if err != nil {
			return nil, fmt.Errorf("error of parse item(%v) while list items", si.ID)
		}

		l = append(l, *i)
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

func parseItem(i *pclient.Item) (*Item, error) {
	b, err := parseBody(&i.Body)
	if err != nil {
		return nil, fmt.Errorf("error of parse item(%v) body:%w", i.ID, err)
	}

	return &Item{
		ID:         i.ID,
		Body:       b,
		Meta:       Meta(i.Body.Meta),
		CreateTime: i.CreateTime,
		UpdateTime: i.UpdateTime,
	}, nil
}

func convertPassword(b *Password) *model.Password {
	return &model.Password{
		Resource: b.Resource,
		UserName: b.UserName,
		Password: b.Password,
	}
}

func convertCard(b *Card) *model.Card {
	return &model.Card{
		Number:  b.Number,
		Expires: b.Expires,
		Holder:  b.Holder,
	}
}

func convertNote(b *Note) *model.Note {
	return &model.Note{
		Name: b.Name,
		Body: b.Body,
	}
}

func convertFile(b *File) *model.File {
	return &model.File{
		Name: b.Name,
		Body: b.Body,
	}
}

func convertAndFillBody(i *model.Item, body any) error {
	switch b := body.(type) {
	case *Password:
		i.Password = convertPassword(b)
	case *Card:
		i.Card = convertCard(b)
	case *Note:
		i.Note = convertNote(b)
	case *File:
		i.File = convertFile(b)
	default:
		return fmt.Errorf("unkown item body type:%v", reflect.TypeOf(b))
	}

	return nil
}

func parsePassword(b *model.Password) *Password {
	return &Password{
		Resource: b.Resource,
		UserName: b.UserName,
		Password: b.Password,
	}
}

func parseCard(b *model.Card) *Card {
	return &Card{
		Number:  b.Number,
		Expires: b.Expires,
		Holder:  b.Holder,
	}
}

func parseNote(b *model.Note) *Note {
	return &Note{
		Name: b.Name,
		Body: b.Body,
	}
}

func parseFile(b *model.File) *File {
	return &File{
		Name: b.Name,
		Body: b.Body,
	}
}

func parseBody(i *model.Item) (any, error) {
	ib, err := i.GetBody()
	if err != nil {
		return nil, fmt.Errorf("error of get item body:%w", err)
	}

	var pib any

	switch b := ib.(type) {
	case *model.Password:
		pib = parsePassword(b)
	case *model.Card:
		pib = parseCard(b)
	case *model.Note:
		pib = parseNote(b)
	case *model.File:
		pib = parseFile(b)
	default:
		return nil, fmt.Errorf("unkown storage item body:%v", reflect.TypeOf(b))
	}

	return pib, nil
}
