package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto/v1"
	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
	"github.com/rs/zerolog/log"
)

type Client interface {
	UserAuthentication
	AuthTokenGeter
	ItemManager
}

type AuthTokenGeter interface {
	GetAuthToken() string
}

type UserAuthentication interface {
	LoginUser(ctx context.Context, login, password string) error
	RegisterUser(ctx context.Context, login, password string) error
	Logout(ctx context.Context)
}

type ItemManager interface {
	GetItem(ctx context.Context, id int64) (*Item, error)
	ListItems(ctx context.Context) ([]Item, error)
	CreateItem(ctx context.Context, item *Item) (int64, error)
	UpdateItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, id int64) error
}

type client struct {
	usersService   pb.UsersServiceClient
	itemsService   pb.ItemsServiceClient
	requestTimeout time.Duration
	authToken      string
}

// New – создание клиента.
func New(a string, rt time.Duration) (*client, error) {
	log.Printf("New grpc client, server address:%v, request timeout:%v seconds", a, rt.Seconds())

	c := &client{
		requestTimeout: rt,
	}

	cc, err := grpc.NewClient(
		a,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AddAuthToken(c)),
	)
	if err != nil {
		return nil, fmt.Errorf("create client error:%w", err)
	}

	c.usersService = pb.NewUsersServiceClient(cc)
	c.itemsService = pb.NewItemsServiceClient(cc)

	return c, nil
}

// Login – логин пользователя на сервере, получение токена.
func (c *client) LoginUser(ctx context.Context, login, password string) error {
	log.Ctx(ctx).Printf("LoginUser, Login:%s", login)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.LoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := c.usersService.Login(ctx, req)
	if err != nil {
		return fmt.Errorf("users service login error:%w", err)
	}

	log.Ctx(ctx).Printf("LoginUser success")
	c.setAuthToken(resp.Token)

	return nil
}

// Logout – логаут пользователя.
func (c *client) Logout(ctx context.Context) {
	log.Ctx(ctx).Printf("Logout => erase auth token")
	c.setAuthToken("")
}

// Register – регистрация пользователя на сервере.
func (c *client) RegisterUser(ctx context.Context, login, password string) error {
	log.Ctx(ctx).Printf("RegisterUser, Login:%s", login)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}
	_, err := c.usersService.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("users service register error:%w", err)
	}

	log.Ctx(ctx).Printf("RegisterUser success")
	return nil
}

// setAuthToken – метод выставления AuthToken пользователя.
func (c *client) setAuthToken(v string) {
	log.Printf("setAuthToken")
	c.authToken = v
}

// GetAuthToken – метод получения AuthToken пользователя.
func (c *client) GetAuthToken() string {
	log.Printf("GetAuthToken")
	return c.authToken
}

// GetItem – получить предмет пользователя.
func (c *client) GetItem(ctx context.Context, id int64) (*Item, error) {
	log.Ctx(ctx).Printf("GetItem, id:%v", id)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.GetItemRequest{
		Id: id,
	}
	resp, err := c.itemsService.GetItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("items service get error:%w", err)
	}

	body, err := model.Deserialize(resp.Item.Data)
	if err != nil {
		return nil, fmt.Errorf("error of deserialize item while get item:%w", err)
	}

	log.Ctx(ctx).Printf("GetItem success")
	return &Item{
		ID:         resp.Item.Id,
		Body:       *body,
		CreateTime: resp.Item.CreateTime.AsTime(),
		UpdateTime: resp.Item.UpdateTime.AsTime(),
	}, nil
}

// ListItems – получить все предметы.
func (c *client) ListItems(ctx context.Context) ([]Item, error) {
	log.Ctx(ctx).Printf("ListItems")

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.ListItemsRequest{}
	resp, err := c.itemsService.ListItems(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("items service list error:%w", err)
	}

	items := make([]Item, 0, len(resp.Items))
	for _, i := range resp.Items {
		body, err := model.Deserialize(i.Data)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("error of deserialize item(%v) while list items", i.Id)
			continue
		}

		item := Item{
			ID:         i.Id,
			Body:       *body,
			CreateTime: i.CreateTime.AsTime(),
			UpdateTime: i.UpdateTime.AsTime(),
		}
		items = append(items, item)
	}

	return items, nil
}

// CreateItem – создать предмет.
func (c *client) CreateItem(ctx context.Context, item *Item) (int64, error) {
	log.Ctx(ctx).Printf("CreateItem, local id:%v", item.ID)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	b, err := model.Serialize(&item.Body)
	if err != nil {
		return 0, fmt.Errorf("error of serialize item(%v) while create item:%w", item.ID, err)
	}

	req := &pb.CreateItemRequest{
		Item: &pb.Item{
			Data:       b,
			CreateTime: timestamppb.New(item.CreateTime),
			UpdateTime: timestamppb.New(item.UpdateTime),
		},
	}
	resp, err := c.itemsService.CreateItem(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("items client create error:%w", err)
	}

	log.Ctx(ctx).Printf("CreateItem success, remote id:%v", resp.Id)
	return resp.Id, nil
}

// UpdateItem – обновить предмет.
func (c *client) UpdateItem(ctx context.Context, item *Item) error {
	log.Ctx(ctx).Printf("UpdateItem, id:%v", item.ID)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	b, err := model.Serialize(&item.Body)
	if err != nil {
		return fmt.Errorf("error of serialize item(%v) while update item:%w", item.ID, err)
	}

	req := &pb.UpdateItemRequest{
		Item: &pb.Item{
			Id:         item.ID,
			Data:       b,
			CreateTime: timestamppb.New(item.CreateTime),
			UpdateTime: timestamppb.New(item.UpdateTime),
		},
	}
	_, err = c.itemsService.UpdateItem(ctx, req)
	if err != nil {
		return fmt.Errorf("items client update item data error:%w", err)
	}

	return nil
}

// DeleteItem – удалить предмет.
func (c *client) DeleteItem(ctx context.Context, id int64) error {
	log.Ctx(ctx).Printf("DeleteItem, id:%v", id)

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	req := &pb.DeleteItemRequest{
		Id: id,
	}

	_, err := c.itemsService.DeleteItem(ctx, req)
	if err != nil {
		return fmt.Errorf("items client update item data error:%w", err)
	}

	return nil
}
