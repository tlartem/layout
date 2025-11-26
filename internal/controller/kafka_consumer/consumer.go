package kafka_consumer

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"gitlab.noway/internal/usecase"
	"gitlab.noway/pkg/logger"
	"gitlab.noway/pkg/metrics"
	"gitlab.noway/pkg/otel"
	"gitlab.noway/pkg/otel/tracer"
)

type Config struct {
	Addr     []string `envconfig:"KAFKA_CONSUMER_ADDR"     required:"true"`
	Topic    string   `default:"noway-layout-topic"  envconfig:"KAFKA_CONSUMER_TOPIC"`
	Group    string   `default:"noway-layout-group"  envconfig:"KAFKA_CONSUMER_GROUP"`
	Disabled bool     `envconfig:"KAFKA_CONSUMER_DISABLED"`
}

type Consumer struct {
	config  Config
	reader  *kafka.Reader
	usecase *usecase.UseCase
	metrics *metrics.Entity
	stop    context.CancelFunc
	done    chan struct{}
}

func New(cfg Config, m *metrics.Entity, uc *usecase.UseCase) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          cfg.Addr,
		Topic:            cfg.Topic,
		GroupID:          cfg.Group,
		ErrorLogger:      logger.ErrorLogger(),
		ReadBatchTimeout: time.Second,
		// CommitInterval: time.Second,
	})

	ctx, stop := context.WithCancel(context.Background())

	c := &Consumer{
		config:  cfg,
		reader:  r,
		usecase: uc,
		metrics: m,
		stop:    stop,
		done:    make(chan struct{}),
	}

	if c.config.Disabled {
		log.Info().Msg("kafka consumer: disabled")

		return c
	}

	go c.run(ctx)

	return c
}

func (c *Consumer) run(ctx context.Context) {
	const consume = "consume"

	log.Info().Msg("kafka consumer: started")

	for {
		now := time.Now()

		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("kafka consumer: FetchMessage")

			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				break
			}
		}

		ctx := otel.ExtractPropagateHeaders(ctx, m)

		ctx, span := tracer.Start(ctx, "kafka consumer from "+c.config.Topic,
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				semconv.MessagingSystemKafka,
				semconv.MessagingDestinationSubscriptionName(c.config.Topic),
				semconv.MessagingConsumerGroupName(c.config.Group),
				semconv.MessagingKafkaMessageKey(string(m.Key)),
			),
		)

		err = c.usecase.Consume(ctx, m)
		if err != nil {
			c.metrics.Total(consume, metrics.Error)
			log.Error().Err(err).Msg("kafka consumer: some work failed")
		}

		if err = c.reader.CommitMessages(ctx, m); err != nil {
			c.metrics.Total(consume, metrics.Error)
			log.Error().Err(err).Msg("kafka consumer: CommitMessages")
		}

		c.metrics.Duration(consume, now)
		c.metrics.Total(consume, metrics.Ok)

		span.End() // Закрываем span
	}

	close(c.done)
}

func (c *Consumer) Close() {
	if c.config.Disabled {
		return
	}

	log.Info().Msg("kafka consumer: closing")

	c.stop()

	if err := c.reader.Close(); err != nil {
		log.Error().Err(err).Msg("kafka consumer: reader.Close")
	}

	<-c.done

	log.Info().Msg("kafka consumer: closed")
}
