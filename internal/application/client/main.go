// Package client for store private data.
package client

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	grpcclient "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	"github.com/k0st1a/gophkeeper/internal/application/client/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/auth"
	"github.com/k0st1a/gophkeeper/internal/pkg/logwrap"
	"github.com/rs/zerolog/log"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config create error:%w", err)
	}
	log.Printf("Cfg:%+v", cfg)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancelFunc()

	err = logwrap.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logwrap create error:%w", err)
	}

	auth := auth.New(cfg.SecretKey)

	srv, err := grpcclient.New(cfg, db, auth, db)
	if err != nil {
		return fmt.Errorf("make grpc client error:%w", err)
	}

	go func() {
		err := srv.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to run client")
		}
	}()

	<-ctx.Done()

	err = srv.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown client")
	}

	return nil
}
