package otel

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"gitlab.noway/pkg/otel/tracer"
)

func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Создаем корневой span
	ctx, span := tracer.Start(ctx, "grpc "+info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	// Вызываем следующий обработчик (или сам handler)
	resp, err := handler(ctx, req)

	// Записываем полезные атрибуты
	span.SetAttributes(
		semconv.HTTPRoute(info.FullMethod),
	)

	// Помечаем span как ошибочный для 4xx и 5xx статусов
	if err != nil {
		span.SetStatus(codes.Error, "")
		span.AddEvent("error", trace.WithAttributes(
			attribute.String("error.message", err.Error()),
		))
	}

	return resp, err
}
