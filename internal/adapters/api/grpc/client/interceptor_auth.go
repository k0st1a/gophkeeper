package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthTokenGeter interface {
	GetAuthToken() string
}

// AddAuthToken – интерсептор, который добавляет token аутентикации в запрос.
func AddAuthToken(g AuthTokenGeter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", g.GetAuthToken())
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
