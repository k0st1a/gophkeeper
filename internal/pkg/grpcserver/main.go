// Package grpcserver is some behaviour of GRPC server.
package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Listener *net.Listener
	Server   *grpc.Server
}

// New create grpc server, where:
//   - address - server host and port;
//   - server - launcing grpc api.
func New(address string, server *grpc.Server) (*Server, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("net listen error:%w", err)
	}

	return &Server{
		Listener: &l,
		Server:   server,
	}, nil
}

// Run - запуск сервера.
func (s *Server) Run() error {
	log.Printf("Run grpc api")

	err := s.Server.Serve(*s.Listener)
	if err != nil {
		return fmt.Errorf("server listen error:%w", err)
	}

	return nil
}

// Shutdown - graceful выключение сервера.
func (s *Server) Shutdown() error {
	s.Server.GracefulStop()
	return nil
}
