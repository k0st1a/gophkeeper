package model

// Note - описание заметки.
//
//easyjson:json
type Note struct {
	Name string `json:"name"` // Название заметки
	Body string `json:"body"` // Тело заметки
}

func (n *Note) GetName() string {
	return n.Name
}
