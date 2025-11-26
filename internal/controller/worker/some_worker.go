package worker

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"gitlab.noway/internal/usecase"
	"gitlab.noway/pkg/otel/tracer"
)

type SomeWorker struct {
	usecase    *usecase.UseCase
	cron       *cron.Cron
	timeToWork chan struct{}
	stop       chan struct{}
	done       chan struct{}
}

func NewSomeWorker(uc *usecase.UseCase) (*SomeWorker, error) {
	w := &SomeWorker{
		usecase:    uc,
		cron:       cron.New(),
		timeToWork: make(chan struct{}),
		stop:       make(chan struct{}),
		done:       make(chan struct{}),
	}

	go w.run()

	// Cron config
	_, err := w.cron.AddFunc("0 8 * * *", func() { // 8:00 every day
		select {
		case w.timeToWork <- struct{}{}:
		default:
		}
	})
	if err != nil {
		return nil, fmt.Errorf("cron.AddFunc: %w", err)
	}

	w.cron.Start()

	return w, nil
}

func (w *SomeWorker) run() {
	log.Info().Msg("some worker: started")
FOR:
	for {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "worker some", trace.WithSpanKind(trace.SpanKindInternal))

		err := w.usecase.SomeWork(ctx)
		if err != nil {
			log.Error().Err(err).Msg("some worker: some work failed")
		}

		span.End() // Закрываем span

		select {
		case <-w.stop:
			break FOR
		case <-w.timeToWork:
		}
	}

	log.Info().Msg("some worker: stopped")

	close(w.done)
}

func (w *SomeWorker) Stop() {
	close(w.stop)

	w.cron.Stop()

	<-w.done
}
