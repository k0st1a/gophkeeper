package model

// Card - описание карты.
//
//easyjson:json
type Card struct {
	Number  string `json:"number"`  // Номер карты
	Expires string `json:"expires"` // Время истечения карты
	Holder  string `json:"holder"`  // Держатель карты
}

func (c *Card) GetName() string {
	return c.Number
}
