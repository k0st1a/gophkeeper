// Package password for serialize/deserialize password model to/from JSON.
package password

import (
	"fmt"

	"github.com/mailru/easyjson"
)

//go:generate easyjson -all model.go

// Password - описание пароля.
//
//easyjson:json
type Password struct {
	UserName string `json:"user_name"` // Имя пользователя
	Password string `json:"password"`  // Пароль пользователя
}

// Deserialize - распаковка байт в формат Password.
func Deserialize(b []byte) (*Password, error) {
	m := &Metrics{}
	err := easyjson.Unmarshal(b, m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Unmarshal error:%w", err)
	}

	return m, nil
}

// Serialize - упаковка Password в байты.
func Serialize(m *Password) ([]byte, error) {
	b, err := easyjson.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("easyjson.Marshal error:%w", err)
	}

	return b, nil
}
