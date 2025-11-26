//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"gitlab.noway/pkg/otel"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.noway/config"
	kafka_producer "gitlab.noway/internal/adapter/kafka"
	"gitlab.noway/internal/app"
	"gitlab.noway/internal/controller/grpc"
	"gitlab.noway/internal/controller/kafka_consumer"
	"gitlab.noway/internal/controller/worker"
	"gitlab.noway/pkg/httpserver"
	"gitlab.noway/pkg/logger"
	"gitlab.noway/pkg/postgres"
	"gitlab.noway/pkg/redis"
)

// Prepare:  make up
// Run test: make integration-test

var ctx = context.Background()

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions

	profile     *ProfileClient
	kafkaWriter *kafka.Writer
	kafkaReader *kafka.Reader
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	s.ResetMigrations()

	// Config
	c := config.Config{
		App: config.App{
			Name:    "layout",
			Version: "test",
		},
		HTTP: httpserver.Config{
			Port: "8080",
		},
		GRPC: grpc.Config{
			Port: "50051",
		},
		Logger: logger.Config{
			AppName:       "layout",
			AppVersion:    "test",
			Level:         "debug",
			PrettyConsole: true,
		},
		Postgres: postgres.Config{
			Host:     "localhost",
			Port:     "5432",
			User:     "login",
			Password: "pass",
			DBName:   "postgres",
		},
		Redis: redis.Config{
			Addr: "localhost:6379",
		},
		ProduceWorker: worker.ProduceConfig{
			Timeout:      time.Second,
			MessageCount: 1,
		},
		KafkaProducer: kafka_producer.Config{
			Addr:  []string{"localhost:9094"},
			Topic: "",
		},
		KafkaConsumer: kafka_consumer.Config{
			Addr:     []string{"localhost:9094"},
			Topic:    "noway-layout-topic",
			Group:    "noway-layout-group",
			Disabled: true, // Disable consumer in test!
		},
		OutboxKafkaWorker: worker.OutboxKafkaConfig{
			Limit: 10,
		},
	}

	// logger.Init(c.Logger)
	log.Logger = zerolog.Nop()
	otel.SilentModeInit()

	// Kafka writer for direct produce messages
	s.kafkaWriter = &kafka.Writer{
		Addr:  kafka.TCP(c.KafkaProducer.Addr...),
		Topic: c.KafkaProducer.Topic,
	}

	// Kafka reader for direct consume messages
	s.kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.KafkaConsumer.Addr,
		Topic:   c.KafkaConsumer.Topic,
		GroupID: c.KafkaConsumer.Group,
	})

	// Server
	go func() {
		err := app.Run(context.Background(), c)
		s.NoError(err)
	}()

	BuildProfile(s)

	time.Sleep(time.Second)
}

func (s *Suite) TearDownSuite() {}

func (s *Suite) SetupTest() {}

func (s *Suite) TearDownTest() {}
