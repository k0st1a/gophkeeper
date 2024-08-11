package client

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcClient – grpc клиента для взаимодействия с grpc сервером.
type grpcClient struct {
	usersClient    pb.UsersServiceClient
	itemsClient    pb.ItemsServiceClient
	authToken      string
	serverAddress  string
	requestTimeout time.Duration
}

// New – функция инициализации gRPC клиента.
func New(serverAddress string, requestTimeout time.Duration) (*grpcClient, error) {
	client := &Client{
		config:  c,
		timeout: time.Duration(c.ConnectionTimeout) * time.Second,
	}

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AddAuthToken(client)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc connection refused: %w", err)
	}

	return &Client{
		usersClient: pb.NewUsersServiceClient(conn),
		itemsClient: pb.NewItemsServiceClient(conn),
		serverAddress:  serverAddress,
		requestTimeout: requestTimeout,
	}, nil
}

// SetAuthToken – метод выставления AuthToken пользователя.
func (c *grpcClient) SetAuthToken(v string) {
	c.authToken = v
	return
}

// GetAuthToken – метод получения AuthToken пользователя.
func (c *grpcClient) GetAuthToken() string {
	return c.authToken
}

// Login – логин пользователя на сервере, получение токена и его сохранение.
func (c *grpcClient) LoginUser(ctx context.Context, login, password string) error {
	log.Ctx(ctx).Printf("LoginUser, Login:%s", login)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := c.usersClient.Login(ctx, req)
	if err != nil {
		return fmt.Errorf("users client login error:%w", err)
	}

	c.SetAuthToken(resp.Token)
	log.Ctx(ctx).Printf("LoginUser success, Token:%v", resp.Token)
	return nil
}

// Register – регистрация пользователя на сервере.
func (c *grpcClient) RegisterUser(ctx context.Context, login, password string) error {
	log.Ctx(ctx).Printf("RegisterUser, Login:%s", login)

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}
	_, err := c.usersClient.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("users client register error:%w", err)
	}

	log.Ctx(ctx).Printf("RegisterUser success")
	return nil
}

// GetItem – получение итема пользователя с сервера.
func (c *grpcClient) GetItem(itemID int64) (Item, error) {
	log.Ctx(ctx).Printf("GetItem, itemID:%v", itemID)

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	req := &pb.GetRequest{
		Id: itemID,
	}
	resp, err := c.itemsClient.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("items client get error:%w", err)
	}

	log.Ctx(ctx).Printf("GetItem success")
	return &Item{
		ID:   resp.Data.Id,
		Name: resp.Data.Item.Name,
		Type: resp.Data.Item.Type,
		Data: resp.Data.Item.Data,
	}, nil
}

// ListItem – получение всех мета-данных пользователя с сервера.
func (c *grpcClient) ListItem(ctx context.Context) ([]ListItem, error) {
	log.Ctx(ctx).Printf("ListItem")

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.ListRequest{}
	resp, err := c.itemsClient.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("items client list error:%w", err)
	}

	items := make([]ListItem, 0, len(res.Data))
	for _, i := range res.Data {
		item := ListItem{
			ID:       item.Id,
			Name:     item.Name,
			DataType: item.DataType,
		}
		items = append(items, item)
	}

	return items, nil
}

// CreateItem – создать итем пользователя на сервере.
func (c *grpcClient) CreateItem(ctx context.Context, name, dataType string, data []byte) error {
	log.Ctx(ctx).Printf("CreateItem, name:%v, dataType:%v", name, dataType)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.CreateRequest{
		Item: &pb.Item{
			Name: name,
			Type: dataType,
			Data: data,
		}
	}
	_, err := c.gRPCClient.CreateItem(ctx, req)
	if err != nil {
		return fmt.Errorf("items client create error:%w", err)
	}

	log.Ctx(ctx).Printf("CreateItem success")
	return nil
}

// UpdateItemData – обновление данных итема на сервер.
func (c *grpcClient) UpdateItemData(ctx, itemID int64, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.UpdateItemDataRequest{
		Id:   itemID,
		Data: data,
	}
	err := c.itemsClient.UpdateItemData(ctx, req)
	if err != nil {
		return fmt.Errorf("items client update item data error:%w", err)
	}

	return nil
}
