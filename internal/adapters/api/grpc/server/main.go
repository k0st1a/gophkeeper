// Package server for process request by grpc.
package server

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/server/handler"
	"github.com/k0st1a/gophkeeper/internal/application/server/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/pkg/grpcserver"
	"github.com/k0st1a/gophkeeper/internal/pkg/user"
)

func New(cfg *config.Config, u user.Managment, a auth.UserAuthentication) (*grpcserver.Server, error) {
	h := &handler.AuthServer{
		User: u,
		Auth: a,
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	pb.RegisterAuthServiceServer(s, h)

	srv, err := grpcserver.New(cfg.Address, s)
	if err != nil {
		return nil, fmt.Errorf("grpc server new error:%w", err)
	}

	return srv, nil
}
