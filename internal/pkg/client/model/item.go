package model

//go:generate easyjson -all item.go

import (
	"errors"
	"fmt"

	"github.com/mailru/easyjson"
)

var (
	ErrBadItem = errors.New("bad item")
)

// Item - описание предмета клиента.
// Должно быть заполнено одно из полей: Card, Password, Note, File.
//
//easyjson:json
type Item struct {
	// Поле Card заполняется, если предмет содержит информацию о банковской карте.
	Card *Card `json:"card"`
	// Поле Password заполняется, если предмет содержит информацию о пароле.
	Password *Password `json:"password"`
	// Поле Note заполняется, если предмет содержит информацию о заметке (текстовые данные).
	Note *Note `json:"note"`
	// Поле File заполняется, если предмет содержит информацию о файле (бинарные данные).
	File *File `json:"file"`
	// Поле Meta содержит опциональную информацию о предмете.
	Meta Meta `json:"meta"`
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

	return "", ErrBadItem
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

	return "", ErrBadItem
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
