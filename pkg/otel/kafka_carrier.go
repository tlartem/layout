package otel

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

func InjectPropagateHeaders(ctx context.Context, msgs ...kafka.Message) {
	for i, msg := range msgs {
		// Создаем carrier на основе заголовков
		carrier := KafkaHeadersCarrier(msg.Headers)

		// Инжектим trace context в заголовки Кафки
		otel.GetTextMapPropagator().Inject(ctx, &carrier)

		msg.Headers = carrier
		msgs[i] = msg
	}
}

func ExtractPropagateHeaders(ctx context.Context, msg kafka.Message) context.Context {
	// Создаем carrier на основе заголовков
	carrier := KafkaHeadersCarrier(msg.Headers)

	// Извлекаем контекст из headers
	ctx = otel.GetTextMapPropagator().Extract(ctx, &carrier)

	return ctx
}

type KafkaHeadersCarrier []kafka.Header

func (c *KafkaHeadersCarrier) Get(key string) string {
	for _, h := range *c {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}

func (c *KafkaHeadersCarrier) Set(key, value string) {
	*c = append(*c, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *KafkaHeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(*c))

	for _, h := range *c {
		keys = append(keys, h.Key)
	}

	return keys
}
