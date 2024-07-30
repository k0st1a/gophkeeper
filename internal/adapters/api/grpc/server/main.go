// Package server for process request by grpc.
package server

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/protobuf/auth"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/server/handler"
	"github.com/k0st1a/gophkeeper/internal/application/server/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/grpcserver"
)

func New(cfg *config.Config) (*grpcserver.Server, error) {
	h := &handler.AuthServer{}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()

	// регистрируем сервис
	auth.RegisterAuthServer(s, h)

	srv, err := grpcserver.New(cfg.Address, s)
	if err != nil {
		return nil, fmt.Errorf("grpc server new error:%w", err)
	}

	return srv, nil
}
