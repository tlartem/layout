package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"gitlab.noway/internal/adapter/kafka"
	"gitlab.noway/internal/controller/grpc"
	"gitlab.noway/internal/controller/kafka_consumer"
	"gitlab.noway/internal/controller/worker"
	"gitlab.noway/pkg/httpclient"
	"gitlab.noway/pkg/httpserver"
	"gitlab.noway/pkg/logger"
	"gitlab.noway/pkg/otel"
	"gitlab.noway/pkg/postgres"
	"gitlab.noway/pkg/redis"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}

type Config struct {
	App               App
	HTTP              httpserver.Config
	GRPC              grpc.Config
	Logger            logger.Config
	OTEL              otel.Config
	Postgres          postgres.Config
	Redis             redis.Config
	Client            httpclient.Config
	KafkaConsumer     kafka_consumer.Config
	KafkaProducer     kafka.Config
	ProduceWorker     worker.ProduceConfig
	OutboxKafkaWorker worker.OutboxKafkaConfig
}

func New() (Config, error) {
	var config Config

	err := godotenv.Load(".env")
	if err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	err = envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
