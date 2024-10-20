package model

//go:generate easyjson -all item.go

import (
	"errors"
	"fmt"

	"github.com/mailru/easyjson"
)

var (
	ErrorBadItem = errors.New("bad item")
)

// Item - описание предмета клиента.
//
//easyjson:json
type Item struct {
	Card     *Card     `json:"card"`
	Password *Password `json:"password"`
	Note     *Note     `json:"note"`
	File     *File     `json:"file"`
}

func (i *Item) GetBody() (any, error) {
	if i.Card != nil {
		return i.Card, nil
	}

	if i.Password != nil {
		return i.Password, nil
	}

	if i.Note != nil {
		return i.Note, nil
	}

	if i.File != nil {
		return i.File, nil
	}

	return "", ErrorBadItem
}

func (i *Item) GetName() (string, error) {
	if i.Card != nil {
		return i.Card.GetName(), nil
	}

	if i.Password != nil {
		return i.Password.GetName(), nil
	}

	if i.Note != nil {
		return i.Note.GetName(), nil
	}

	if i.File != nil {
		return i.File.GetName(), nil
	}

	return "", ErrorBadItem
}

// Deserialize - распаковка байт в формат Item.
func Deserialize(b []byte) (*Item, error) {
	i := &Item{}
	err := easyjson.Unmarshal(b, i)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	return i, nil
}

// Serialize - упаковка Item в байты.
func Serialize(i *Item) ([]byte, error) {
	b, err := easyjson.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}
