package model

// Password - описание пароля.
//
//easyjson:json
type Password struct {
	Resource string `json:"resource"`  // Пароль для данного ресурса
	UserName string `json:"user_name"` // Имя пользователя
	Password string `json:"password"`  // Пароль пользователя
}

func (p *Password) GetName() string {
	return p.Resource
}
