package otel

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/router"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем контекст из заголовков запроса
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Создаем корневой span
		ctx, span := tracer.Start(ctx, "", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// Оборачиваем writer для захвата статуса ответа
		ww := router.WriterWrapper(w)

		// Вызываем следующий обработчик (или сам handler)
		next.ServeHTTP(ww, r.WithContext(ctx))

		span.SetName("http " + r.Method + " " + router.ExtractPath(ctx))

		// Записываем полезные атрибуты
		span.SetAttributes(
			semconv.HTTPResponseStatusCode(ww.Code()),
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPRoute(r.URL.Path),
		)

		// Помечаем span как ошибочный для 4xx и 5xx статусов
		if ww.Code() >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(ww.Code()))
			span.AddEvent("error", trace.WithAttributes(
				attribute.String("error.message", http.StatusText(ww.Code())),
			))
		}
	}

	return http.HandlerFunc(fn)
}
