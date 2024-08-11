// Package client for store private data.
package client

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	cliclient "github.com/k0st1a/gophkeeper/internal/adapters/api/cli/client"
	grpcclient "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	"github.com/k0st1a/gophkeeper/internal/application/client/config"
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

	grpc, err := grpcclient.New(cfg.Address, 3*time.Second)
	if err != nil {
		return fmt.Errorf("make grpc client error:%w", err)
	}

	cli, err := cliclient.New(grpc)
	if err != nil {
		return fmt.Errorf("make cli client error:%w", err)
	}

	go func() {
		err := cli.Run()
		if err != nil {
			log.Error().Err(err).Msg("failed to run cli client")
		}
	}()

	<-ctx.Done()

	err = cli.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("error of shutdown cli client")
	}

	return nil
}
