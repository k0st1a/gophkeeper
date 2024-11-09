// Package server for process request by grpc.
package server

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto/v1"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/server/handler"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/server/interceptor"
	"github.com/k0st1a/gophkeeper/internal/application/server/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/pkg/grpcserver"
	"github.com/k0st1a/gophkeeper/internal/ports/server"
)

func New(cfg *config.Config, u server.UserStorage, a auth.UserAuthentication,
	i server.ItemStorage) (*grpcserver.Server, error) {
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.Authenticate(a),
	))

	uh := &handler.UserServer{
		Storage: u,
		Auth:    a,
	}

	pb.RegisterUsersServiceServer(s, uh)

	ih := &handler.ItemServer{
		Storage: i,
	}
	pb.RegisterItemsServiceServer(s, ih)

	srv, err := grpcserver.New(cfg.Address, s)
	if err != nil {
		return nil, fmt.Errorf("grpc server new error:%w", err)
	}

	return srv, nil
}
