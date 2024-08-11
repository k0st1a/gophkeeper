package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/pkg/userid"
	"github.com/rs/zerolog/log"
)

func Authenticate(auth auth.UserAuthentication) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, r interface{}, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		if i.FullMethod == pb.UsersService_Register_FullMethodName ||
			i.FullMethod == pb.UsersService_Login_FullMethodName {
			return h(ctx, r)
		}

		var token string
		if meta, ok := metadata.FromIncomingContext(ctx); ok {
			values := meta.Get("token")
			if len(values) > 0 {
				token = values[0]
			}
		}

		if len(token) == 0 {
			log.Ctx(ctx).Printf("no token")
			return nil, status.Errorf(codes.Unauthenticated, "no token")
		}

		userID, err := auth.GetUserID(token)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("error of get userID")
			return nil, status.Errorf(codes.Unauthenticated, "no user id in token")
		}

		CtxWithUserID := userid.Add(ctx, userID)

		return h(CtxWithUserID, r)
	}
}
