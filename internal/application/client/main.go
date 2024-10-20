// Package client for store private data.
package client

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/tui"
	tstorage "github.com/k0st1a/gophkeeper/internal/adapters/api/tui/storage"
	"github.com/k0st1a/gophkeeper/internal/adapters/storage/inmemory"
	"github.com/k0st1a/gophkeeper/internal/application/client/config"
	"github.com/k0st1a/gophkeeper/internal/pkg/job"
	"github.com/k0st1a/gophkeeper/internal/pkg/logwrap"
	itemsync "github.com/k0st1a/gophkeeper/internal/pkg/sync"
	"github.com/k0st1a/gophkeeper/internal/pkg/tick"
	"github.com/rs/zerolog/log"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config create error:%w", err)
	}
	log.Printf("Cfg:%+v", cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if cfg.LogFile != "" {
		err = logwrap.NewFile(cfg.LogLevel, cfg.LogFile)
	} else {
		err = logwrap.New(cfg.LogLevel)
	}

	if err != nil {
		return fmt.Errorf("logwrap create error:%w", err)
	}

	gc, err := client.New(cfg.Address, 3*time.Second)
	if err != nil {
		return fmt.Errorf("make grpc client error:%w", err)
	}

	s := inmemory.New()

	is := itemsync.New(s, gc)
	t := tick.New(is, 10*time.Second)

	j := job.New(t)

	ctx, cancel = context.WithCancel(ctx)

	ts := tstorage.New(s)

	ui := tui.New(gc, ts, j, cancel)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := ui.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to run ui")
		}
	}()

	<-ctx.Done()

	ui.Stop(ctx)

	wg.Wait()

	return nil
}
