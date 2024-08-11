package client

type ListItem struct {
	ID       int64
	Name     string
	DataType string
}

type Item struct {
	ID       int64
	Name     string
	DataType string
	Data     []byte
}
