// Package handler contains handler for grpc server.
package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/ports"
	"github.com/rs/zerolog/log"
)

type UserServer struct {
	// нужно встраивать тип auth.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedUsersServiceServer
	Storage ports.UserStorage
	Auth    auth.UserAuthentication
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Ctx(ctx).Printf("Register, Login:%s", req.Login)

	passwordHash, err := s.Auth.GeneratePasswordHash(req.Password)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("error of generate password pash")
		return nil, status.Errorf(codes.Internal, "create user error")
	}

	id, err := s.Storage.CreateUser(ctx, req.Login, passwordHash)
	if err != nil {
		if errors.Is(err, ports.ErrLoginAlreadyBusy) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		log.Error().Err(err).Ctx(ctx).Msg("error of create user")
		return nil, status.Errorf(codes.Internal, "create user error")
	}

	resp := pb.RegisterResponse{}

	log.Ctx(ctx).Printf("Success register, Login:%s, UserId:%d", req.Login, id)
	return &resp, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Ctx(ctx).Printf("Login, Login:%s", req.Login)

	userID, password, err := s.Storage.GetUserIDAndPassword(ctx, req.Login)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
		}

		log.Error().Err(err).Ctx(ctx).Msg("error of get user id and password")
		return nil, status.Errorf(codes.Internal, "login user error")
	}

	err = s.Auth.CheckPasswordHash(req.Password, password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	t, err := s.Auth.GenerateToken(userID)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("error of generate token")
		return nil, status.Errorf(codes.Internal, "login user error")
	}

	resp := pb.LoginResponse{
		Token: t,
	}

	log.Ctx(ctx).Printf("Success login, Login:%s, UserId:%d", req.Login, userID)
	return &resp, nil
}
