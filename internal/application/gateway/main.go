package gateway

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/gen/proto"
	"github.com/k0st1a/gophkeeper/third_party"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
//
//nolint:lll // not need there
func getOpenAPIHandler() (http.Handler, error) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return nil, fmt.Errorf("mime add extension type error:%w", err)
	}

	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		return nil, fmt.Errorf("couldn't create sub filesystem: %w", err)
	}

	return http.FileServer(http.FS(subFS)), nil
}

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(ctx context.Context, dialAddr string, gatewayAddr string) error {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.NewClient(dialAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create grpc client error:%w", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterUsersServiceHandler(ctx, gwmux, conn)
	if err != nil {
		return fmt.Errorf("failed to register users gateway: %w", err)
	}

	err = pb.RegisterItemsServiceHandler(ctx, gwmux, conn)
	if err != nil {
		return fmt.Errorf("failed to register items gateway: %w", err)
	}

	oa, err := getOpenAPIHandler()
	if err != nil {
		return err
	}

	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				gwmux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}

	log.Info().Msgf("Serving gRPC-Gateway and OpenAPI Documentation on http://%s", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}
