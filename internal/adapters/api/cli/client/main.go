package client

import (
	grpc "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
)

// CLI – взаимодействие с пользователем по средствам cli.
type CLI struct {
	grpc grpc.GrpcClient
}

// New – создание cli клиента.
func New(c grpc.GrpcClient) (*CLI, error) {
	client := CLI{
		grpc: c,
	}

	return &client, nil
}

func (c *CLI) Run() error {
	return nil
}

func (c *CLI) Shutdown() error {
	return nil
}
