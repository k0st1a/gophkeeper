package job

import (
	"context"
	"sync"

	"github.com/k0st1a/gophkeeper/internal/pkg/tick"
	"github.com/rs/zerolog/log"
)

type StartStopper interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type job struct {
	cancel func()
	do     tick.Runner
	wg     sync.WaitGroup
}

func New(do tick.Runner) *job {
	return &job{
		do:     do,
		cancel: func() {},
	}
}

func (j *job) Start(ctx context.Context) {
	log.Ctx(ctx).Printf("Start job")

	ctx, cancel := context.WithCancel(ctx)
	j.cancel = cancel

	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		err := j.do.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to run Runner")
		}
	}()
	log.Ctx(ctx).Printf("Job started")
}

func (j *job) Stop(ctx context.Context) {
	log.Ctx(ctx).Printf("Stop job")
	j.cancel()
	j.wg.Wait()
	log.Ctx(ctx).Printf("Job stopped")
}
