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
