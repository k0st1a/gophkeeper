package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/rs/zerolog/log"
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

type GrpcClient interface {
	UserAuthentication
	ItemManager
}

type UserAuthentication interface {
	LoginUser(ctx context.Context, login, password string) error
	RegisterUser(ctx context.Context, login, password string) error
}

type ItemManager interface {
	GetItem(ctx context.Context, itemID int64) (*Item, error)
	ListItem(ctx context.Context) ([]ListItem, error)
	CreateItem(ctx context.Context, name, dataType string, data []byte) error
	UpdateItemData(ctx context.Context, itemID int64, data []byte) error
}

// New – функция инициализации gRPC клиента.
func New(serverAddress string, requestTimeout time.Duration) (*grpcClient, error) {
	log.Printf("New grpc client, serverAddress:%v, requestTimeout:%v", serverAddress, requestTimeout.Seconds())

	client := &grpcClient{
		serverAddress:  serverAddress,
		requestTimeout: requestTimeout,
	}
	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AddAuthToken(client)),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc connection refused: %w", err)
	}

	client.usersClient = pb.NewUsersServiceClient(conn)
	client.itemsClient = pb.NewItemsServiceClient(conn)

	return client, nil
}

// SetAuthToken – метод выставления AuthToken пользователя.
func (c *grpcClient) SetAuthToken(v string) {
	log.Printf("SetAuthToken, token:%v", v)
	c.authToken = v
	return
}

// GetAuthToken – метод получения AuthToken пользователя.
func (c *grpcClient) GetAuthToken() string {
	log.Printf("GetAuthToken")
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

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
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
func (c *grpcClient) GetItem(ctx context.Context, itemID int64) (*Item, error) {
	log.Ctx(ctx).Printf("GetItem, itemID:%v", itemID)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
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

	items := make([]ListItem, 0, len(resp.Data))
	for _, i := range resp.Data {
		item := ListItem{
			ID:   i.Id,
			Name: i.Name,
			Type: i.Type,
		}
		items = append(items, item)
	}

	return items, nil
}

// CreateItem – создать итем пользователя на сервере.
func (c *grpcClient) CreateItem(ctx context.Context, name, dataType string, data []byte) error {
	log.Ctx(ctx).Printf("CreateItem, name:%v, dataType:%v", name, dataType)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.CreateRequest{
		Item: &pb.Item{
			Name: name,
			Type: dataType,
			Data: data,
		},
	}
	_, err := c.itemsClient.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("items client create error:%w", err)
	}

	log.Ctx(ctx).Printf("CreateItem success")
	return nil
}

// UpdateItemData – обновление данных итема на сервер.
func (c *grpcClient) UpdateItemData(ctx context.Context, itemID int64, data []byte) error {
	log.Ctx(ctx).Printf("UpdateItemData, itemID:%v", itemID)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.UpdateItemDataRequest{
		Id:   itemID,
		Data: data,
	}
	_, err := c.itemsClient.UpdateItemData(ctx, req)
	if err != nil {
		return fmt.Errorf("items client update item data error:%w", err)
	}

	return nil
}
