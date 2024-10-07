package rawitem

import (
	"bytes"
	"errors"
	"time"
)

var (
	ErrorItemNotFound      = errors.New("item not found")
	ErrorItemAlreadyExists = errors.New("item already exists")
)

// Info - структура для хранения "сырых" предметов.
type Info struct {
	// ID - идентификатор предмета.
	ID string
	// Type - тип предмета
	Type string
	// Name - имя предмета.
	Name string
	// Description - описание предмета.
	Description string
	// Data - байты данных предмета.
	Data []byte
	// CreateTime - время создания предмета.
	CreateTime time.Time
	// UpdateTime - время обновления предмета, выставляется
	// когда предмет был создан, обновлен, удален.
	UpdateTime time.Time
	// Deleted - флаг указывает но то, был ли удален файл.
	Deleted bool
	// UploadTime - Время загрузки предмета на сервер.
	UploadTime time.Time
}
