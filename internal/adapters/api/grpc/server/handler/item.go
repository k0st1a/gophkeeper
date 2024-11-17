// Package handler contains handler for grpc server.
package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto/v1"
	"github.com/k0st1a/gophkeeper/internal/pkg/userid"
	"github.com/k0st1a/gophkeeper/internal/ports/server"
	"github.com/rs/zerolog/log"
)

type ItemServer struct {
	// нужно встраивать тип auth.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedItemsServiceServer
	Storage server.ItemStorage // YAGNI - без промежуточного сервиса логики над item.
}

func (s *ItemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	log.Ctx(ctx).Printf("Create item, CreateTime:%v", req.Item.CreateTime)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Print(ErrNoUserID.Error())
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Unauthenticated, ErrNoUserID.Error())
	}

	item := &server.Item{
		Data:       req.Item.Data,
		CreateTime: req.Item.CreateTime.AsTime(),
		UpdateTime: req.Item.UpdateTime.AsTime(),
	}
	id, err := s.Storage.CreateItem(ctx, userID, item)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("create item error")
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Internal, "create item error")
	}

	resp := pb.CreateItemResponse{
		Id: id,
	}

	log.Ctx(ctx).Printf("Create item success")
	return &resp, nil
}

func (s *ItemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	log.Ctx(ctx).Printf("Update item, id:%v", req.Item.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Print(ErrNoUserID.Error())
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Unauthenticated, ErrNoUserID.Error())
	}

	item := &server.Item{
		ID:         req.Item.Id,
		Data:       req.Item.Data,
		CreateTime: req.Item.CreateTime.AsTime(),
		UpdateTime: req.Item.UpdateTime.AsTime()}
	err := s.Storage.UpdateItem(ctx, userID, item)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("update item error")
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Internal, "update item error")
	}

	log.Ctx(ctx).Printf("Update item data success")
	return &pb.UpdateItemResponse{}, nil
}

func (s *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	log.Ctx(ctx).Printf("Get item, id:%v", req.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Print(ErrNoUserID.Error())
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Unauthenticated, ErrNoUserID.Error())
	}

	i, err := s.Storage.GetItem(ctx, userID, req.Id)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("get item error")
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Internal, "get item error")
	}

	resp := pb.GetItemResponse{
		Item: &pb.Item{
			Id:         i.ID,
			Data:       i.Data,
			CreateTime: timestamppb.New(i.CreateTime),
			UpdateTime: timestamppb.New(i.UpdateTime),
		},
	}

	log.Ctx(ctx).Printf("Get item success")
	return &resp, nil
}

func (s *ItemServer) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	log.Ctx(ctx).Printf("List item")

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Print(ErrNoUserID.Error())
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Unauthenticated, ErrNoUserID.Error())
	}

	l, err := s.Storage.ListItems(ctx, userID)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("list item error")
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Internal, "list item error")
	}

	items := make([]*pb.Item, 0, len(l))
	for _, i := range l {
		d := &pb.Item{
			Id:         i.ID,
			Data:       i.Data,
			CreateTime: timestamppb.New(i.CreateTime),
			UpdateTime: timestamppb.New(i.UpdateTime),
		}
		items = append(items, d)
	}

	resp := pb.ListItemsResponse{
		Items: items,
	}

	log.Ctx(ctx).Printf("Get item success")
	return &resp, nil
}

func (s *ItemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	log.Ctx(ctx).Printf("Delete item, id:%v", req.Id)

	userID, ok := userid.Get(ctx)
	if !ok {
		log.Ctx(ctx).Print(ErrNoUserID.Error())
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Unauthenticated, ErrNoUserID.Error())
	}

	err := s.Storage.DeleteItem(ctx, userID, req.Id)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("delete item error")
		//nolint:wrapcheck // not need wrap error from status package
		return nil, status.Error(codes.Internal, "delte item error")
	}

	log.Ctx(ctx).Printf("Delete item success")
	return &pb.DeleteItemResponse{}, nil
}
