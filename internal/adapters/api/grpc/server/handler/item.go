// Package handler contains handler for grpc server.
package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto/v1"
	"github.com/k0st1a/gophkeeper/internal/pkg/userid"
	"github.com/k0st1a/gophkeeper/internal/ports"
	"github.com/rs/zerolog/log"
)

type ItemServer struct {
	// нужно встраивать тип auth.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedItemsServiceServer
	Storage ports.ItemStorage // YAGNI - без промежуточного сервиса логики над item.
}

func (s *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	log.Ctx(ctx).Printf("Create item, Name:%v, Type:%v", req.Item.Name, req.Item.Type)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Printf("no user id")
		return nil, status.Errorf(codes.Unauthenticated, "no user id")
	}

	id, err := s.Storage.CreateItem(ctx, userID, req.Item.Name, req.Item.Type, req.Item.Data)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("create item error")
		return nil, status.Errorf(codes.Internal, "create item error")
	}

	resp := pb.CreateItemResponse{
		Id: id,
	}

	log.Ctx(ctx).Printf("Create item success")
	return &resp, nil
}

func (s *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	log.Ctx(ctx).Printf("Update item data, itemID:%v", req.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Printf("no user id")
		return nil, status.Errorf(codes.Unauthenticated, "no user id")
	}

	err := s.Storage.UpdateItem(ctx, userID, req.Id, req.Item.Data)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("update item error")
		return nil, status.Errorf(codes.Internal, "update item error")
	}

	log.Ctx(ctx).Printf("Update item data success")
	return &pb.UpdateItemResponse{}, nil
}

func (s *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	log.Ctx(ctx).Printf("Get item, itemID:%v", req.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Printf("no user id")
		return nil, status.Errorf(codes.Unauthenticated, "no user id")
	}

	i, err := s.Storage.GetItem(ctx, userID, req.Id)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("get item error")
		return nil, status.Errorf(codes.Internal, "get item error")
	}

	resp := pb.GetItemResponse{
		ItemInfo: &pb.ItemInfo{
			Id: i.ID,
			Item: &pb.Item{
				Name: i.Name,
				Type: i.Type,
				Data: i.Data,
			},
		},
	}

	log.Ctx(ctx).Printf("Get item success")
	return &resp, nil
}

func (s *ItemServer) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	log.Ctx(ctx).Printf("List item")

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Printf("no user id")
		return nil, status.Errorf(codes.Unauthenticated, "no user id")
	}

	l, err := s.Storage.ListItem(ctx, userID)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("list item error")
		return nil, status.Errorf(codes.Internal, "list item error")
	}

	items := make([]*pb.ItemInfo, 0, len(l))
	for _, i := range l {
		d := &pb.ItemInfo{
			Id: i.ID,
			Item: &pb.Item{
				Name: i.Name,
				Type: i.Type,
			},
		}
		items = append(items, d)
	}

	resp := pb.ListItemsResponse{
		Items: items,
	}

	log.Ctx(ctx).Printf("Get item success")
	return &resp, nil
}
