// Package server for store private data.
package server

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	grpcserver "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/server"
	"github.com/k0st1a/gophkeeper/internal/application/server/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/logwrap"
	"github.com/rs/zerolog/log"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config create error:%w", err)
	}

	err = logwrap.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logwrap create error:%w", err)
	}

	log.Printf("Cfg:%+v", cfg)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelFunc()

	srv, err := grpcserver.New(cfg)
	if err != nil {
		return fmt.Errorf("make grpc server error:%w", err)
	}

	go func() {
		err := srv.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to run server")
		}
	}()

	<-ctx.Done()

	err = srv.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown server")
	}

	return nil
}
