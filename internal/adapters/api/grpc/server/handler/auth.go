// Package handler contains handler for grpc server.
package handler

import (
	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/protobuf/auth"
)

type AuthServer struct {
	// нужно встраивать тип auth.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	auth.UnimplementedAuthServer
}
