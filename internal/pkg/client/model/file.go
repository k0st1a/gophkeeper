package model

// File - описание файла.
//
//easyjson:json
type File struct {
	Name string `json:"name"` // Путь до файла
	Body []byte `json:"body"` // Тело файла
}

func (f *File) GetName() string {
	return f.Name
}
