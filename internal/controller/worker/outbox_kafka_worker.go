package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"gitlab.noway/internal/usecase"
	"gitlab.noway/pkg/otel/tracer"
)

type OutboxKafkaConfig struct {
	Limit int `default:"10" envconfig:"OUTBOX_KAFKA_WORKER_LIMIT"`
}

type OutboxKafkaWorker struct {
	config  OutboxKafkaConfig
	usecase *usecase.UseCase
	stop    chan struct{}
	done    chan struct{}
}

func NewOutboxKafkaWorker(uc *usecase.UseCase, c OutboxKafkaConfig) *OutboxKafkaWorker {
	w := &OutboxKafkaWorker{
		config:  c,
		usecase: uc,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *OutboxKafkaWorker) run() {
	log.Info().Msg("outbox kafka worker: started")

FOR:
	for {
		ctx := context.Background()
		ctx, span := tracer.Start(ctx, "worker outbox kafka", trace.WithSpanKind(trace.SpanKindInternal))

		lenMessages, err := w.usecase.OutboxReadAndProduce(ctx, w.config.Limit)
		if err != nil {
			log.Error().Err(err).Msg("outbox kafka worker: read and produce failed")
		}

		span.End() // Закрываем span

		var sleepDuration time.Duration

		if lenMessages < w.config.Limit {
			sleepDuration = 10 * time.Second
		}

		select {
		case <-w.stop:
			break FOR
		case <-time.After(sleepDuration):
		}
	}

	log.Info().Msg("outbox kafka worker: stopped")

	close(w.done)
}

func (w *OutboxKafkaWorker) Stop() {
	close(w.stop)

	<-w.done
}
