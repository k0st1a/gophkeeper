package storage

import (
	"errors"
	"fmt"
	"time"
)

var (
	MaxFileSize  = 40 * 1024 * 1024 // 40MB
	ErrLargeFile = errors.New("large file")
)

type Item struct {
	ID         string
	Body       any // password, card, file, note
	CreateTime time.Time
	UpdateTime time.Time
}

func (i *Item) GetName() (string, error) {
	switch t := i.Body.(type) {
	case *Password:
		return t.GetName(), nil
	case *Card:
		return t.GetName(), nil
	case *Note:
		return t.GetName(), nil
	case *File:
		return t.GetName(), nil
	}

	return "", fmt.Errorf("unknown item body type")
}

func (i *Item) GetType() (string, error) {
	switch t := i.Body.(type) {
	case *Password:
		return t.GetType(), nil
	case *Card:
		return t.GetType(), nil
	case *Note:
		return t.GetType(), nil
	case *File:
		return t.GetType(), nil
	}

	return "", fmt.Errorf("unknown item body type")
}

type Password struct {
	Resource string
	UserName string
	Password string
}

func (p *Password) GetName() string {
	return p.Resource
}

func (p *Password) GetType() string {
	return "password"
}

type Card struct {
	Number  string
	Expires string
	Holder  string
}

func (c *Card) GetName() string {
	return c.Number
}

func (c *Card) GetType() string {
	return "card"
}

type Note struct {
	Name string
	Body string
}

func (n *Note) GetName() string {
	return n.Name
}

func (n *Note) GetType() string {
	return "note"
}

type File struct {
	Name        string
	Description string
	Body        []byte
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetType() string {
	return "file"
}
