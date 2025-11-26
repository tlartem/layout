package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer //nolint:gochecknoglobals

func Init(t trace.Tracer) {
	tracer = t
}

//nolint:ireturn
func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName, opts...) //nolint:spancheck
}
