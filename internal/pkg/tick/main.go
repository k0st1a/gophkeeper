package tick

import (
	"context"
	"time"

	"github.com/k0st1a/gophkeeper/internal/pkg/sync"
	"github.com/rs/zerolog/log"
)

type Runner interface {
	Run(context.Context) error
}

type tick struct {
	sync     sync.Doer
	interval time.Duration
}

func New(s sync.Doer, i time.Duration) *tick {
	return &tick{
		sync:     s,
		interval: i,
	}
}

func (t *tick) Run(ctx context.Context) error {
	log.Printf("Run ticker, interval:%v seconds", t.interval.Seconds())
	ticker := time.NewTicker(t.interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Ticker closed with cause:%s", ctx.Err())
			return nil
		case <-ticker.C:
			log.Printf("Got tick")
			err := t.sync.Do(ctx)
			if err != nil {
				log.Error().Err(err).Msg("error of do sync")
			}
		}
	}
}
