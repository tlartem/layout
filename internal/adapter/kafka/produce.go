package kafka

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"

	"gitlab.noway/pkg/logger"
	"gitlab.noway/pkg/metrics"
	"gitlab.noway/pkg/otel"
	"gitlab.noway/pkg/otel/tracer"
)

type Config struct {
	Addr  []string `envconfig:"KAFKA_WRITER_ADDR"  required:"true"`
	Topic string   `envconfig:"KAFKA_WRITER_TOPIC"`
}

type Producer struct {
	config  Config
	writer  *kafka.Writer
	metrics *metrics.Entity
}

func NewProducer(c Config, m *metrics.Entity) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		Topic:        c.Topic,
		Balancer:     &kafka.Hash{Hasher: fnv.New32a()},
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  logger.ErrorLogger(),
		// Async:        true,
	}

	return &Producer{
		config:  c,
		writer:  w,
		metrics: m,
	}
}

func (p *Producer) Produce(ctx context.Context, msgs ...kafka.Message) error {
	ctx, span := tracer.Start(ctx, "adapter kafka Produce")
	defer span.End()

	for _, msg := range msgs {
		ctx := otel.ExtractPropagateHeaders(ctx, msg)

		_, span := tracer.Start(ctx, "adapter kafka Produce to "+msg.Topic,
			trace.WithSpanKind(trace.SpanKindProducer))
		defer span.End() //nolint:gocritic
	}

	const produce = "produce"

	defer p.metrics.Duration(produce, time.Now())

	err := p.writer.WriteMessages(ctx, msgs...)
	if err != nil {
		p.metrics.TotalAdd(produce, metrics.Error, len(msgs))

		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	p.metrics.TotalAdd(produce, metrics.Ok, len(msgs))

	return nil
}

func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Error().Err(err).Msg("kafka producer: p.writer.Close")
	}

	log.Info().Msg("kafka producer: closed")
}
