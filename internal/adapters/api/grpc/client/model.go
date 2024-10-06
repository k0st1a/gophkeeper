package client

type ListItem struct {
	ID   int64
	Name string
	Type string
}

type Item struct {
	ID   int64
	Name string
	Type string
	Data []byte
}

type ItemType string

const (
	ItemTypePassword ItemType = "password"
	ItemTypeCard     ItemType = "card"
	ItemTypeNote     ItemType = "note"
	ItemTypeFile     ItemType = "file"
)
