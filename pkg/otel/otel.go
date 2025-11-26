package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	tracer_noop "go.opentelemetry.io/otel/trace/noop"

	"gitlab.noway/pkg/otel/tracer"
)

type Config struct {
	AppName    string  `envconfig:"APP_NAME"`
	AppVersion string  `envconfig:"APP_VERSION"`
	Endpoint   string  `envconfig:"OTEL_ENDPOINT"`
	Namespace  string  `envconfig:"OTEL_NAMESPACE"`
	InstanceID string  `envconfig:"OTEL_INSTANCE_ID"`
	Ratio      float64 `default:"1.0"                envconfig:"OTEL_RATIO"`
}

var shutdownTracing func(ctx context.Context) error //nolint:gochecknoglobals

func SilentModeInit() {
	otel.SetTracerProvider(tracer_noop.NewTracerProvider())
	tracer.Init(otel.Tracer(""))

	log.Info().Msg("otel: Tracer is disabled")
}

func Init(ctx context.Context, c Config) error {
	if c.Endpoint == "" {
		SilentModeInit()

		return nil
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(prop)

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(c.Endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.AppName),
			semconv.ServiceNamespaceKey.String(c.Namespace),
			semconv.ServiceInstanceIDKey.String(c.Namespace),
			semconv.ServiceVersionKey.String(c.Namespace),
		)),
	)

	shutdownTracing = traceProvider.Shutdown

	otel.SetTracerProvider(traceProvider)
	tracer.Init(otel.Tracer(""))

	return nil
}

func Close() {
	if shutdownTracing == nil {
		return
	}

	err := shutdownTracing(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("otel: failed to shutdown tracing")
	}

	log.Info().Msg("otel: closed")
}
