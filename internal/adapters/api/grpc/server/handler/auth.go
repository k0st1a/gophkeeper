// Package handler contains handler for grpc server.
package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/pkg/user"
	"github.com/rs/zerolog/log"
)

type AuthServer struct {
	// нужно встраивать тип auth.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedAuthServiceServer
	User user.Managment
	Auth auth.UserAuthentication
}

func (a *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Ctx(ctx).Printf("Register, Email:%s", req.Email)

	passwordHash, err := a.Auth.GeneratePasswordHash(req.Password)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("error of generate password pash")
		return nil, status.Errorf(codes.Internal, "create user error")
	}

	id, err := a.User.Create(ctx, req.Email, passwordHash)
	if err != nil {
		if errors.Is(err, user.ErrEmailAlreadyBusy) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		log.Error().Err(err).Ctx(ctx).Msg("error of create user")
		return nil, status.Errorf(codes.Internal, "create user error")
	}

	resp := pb.RegisterResponse{
		UserId: id,
	}

	log.Ctx(ctx).Printf("Success register, Email:%s, UserId:%d", req.Email, id)
	return &resp, nil
}

func (a *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Ctx(ctx).Printf("Login, Email:%s", req.Email)

	userID, password, err := a.User.GetIDAndPassword(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
		}

		log.Error().Err(err).Ctx(ctx).Msg("error of get user id and password")
		return nil, status.Errorf(codes.Internal, "login user error")
	}

	err = a.Auth.CheckPasswordHash(req.Password, password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	t, err := a.Auth.GenerateToken(userID)
	if err != nil {
		log.Error().Err(err).Ctx(ctx).Msg("error of generate token")
		return nil, status.Errorf(codes.Internal, "login user error")
	}

	resp := pb.LoginResponse{
		Token: t,
	}

	log.Ctx(ctx).Printf("Success login, Email:%s, UserId:%d", req.Email, userID)
	return &resp, nil
}
