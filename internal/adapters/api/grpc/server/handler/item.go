// Package handler contains handler for grpc server.
package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
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

func (s *ItemServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
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

	resp := pb.CreateResponse{
		Id: id,
	}

	log.Ctx(ctx).Printf("Create item success")
	return &resp, nil
}

func (s *ItemServer) UpdateItemData(ctx context.Context, req *pb.UpdateItemDataRequest) (*pb.UpdateItemDataResponse, error) {
	log.Ctx(ctx).Printf("Update item data, itemID:%v", req.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Printf("no user id")
		return nil, status.Errorf(codes.Unauthenticated, "no user id")
	}

	err := s.Storage.UpdateItem(ctx, userID, req.Id, req.Data)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("update item error")
		return nil, status.Errorf(codes.Internal, "update item error")
	}

	log.Ctx(ctx).Printf("Update item data success")
	return &pb.UpdateItemDataResponse{}, nil
}

func (s *ItemServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
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

	resp := pb.GetResponse{
		Data: &pb.ItemInfo{
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

func (s *ItemServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
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
				Data: i.Data,
			},
		}
		items = append(items, d)
	}

	resp := pb.ListResponse{
		Data: items,
	}

	log.Ctx(ctx).Printf("Get item success")
	return &resp, nil
}
