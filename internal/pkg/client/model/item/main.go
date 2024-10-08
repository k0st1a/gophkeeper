package item

import (
	"errors"
	"time"
)

var (
	ErrorItemNotFound      = errors.New("item not found")
	ErrorItemAlreadyExists = errors.New("item already exists")
)

// Info - структура для хранения информации о предмете на стороне клиента.
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

// Map2List  - преобразование из map в list.
func Map2List(m map[string]Info) []Info {
	l := make([]Info, len(m))
	i := 0
	for _, v := range m {
		l[i] = v
		i++
	}

	return l
}

// Append - добавление в мапу.
func Append(acc map[string]Info, add []Info) map[string]Info {
	for _, v := range add {
		acc[v.Name] = v
	}

	return acc
}
