package usecase

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"gitlab.noway/pkg/otel/tracer"
)

func (u *UseCase) Consume(ctx context.Context, m kafka.Message) error {
	ctx, span := tracer.Start(ctx, "usecase Consume")
	defer span.End()

	if u.redis.IsExists(ctx, string(m.Key)) {
		log.Info().Str("key", string(m.Key)).Msg("usecase: Consume: message already processed")

		return nil
	}

	log.Info().
		Str("topic", m.Topic).
		Int("partition", m.Partition).
		Int64("offset", m.Offset).
		Str("key", string(m.Key)).
		Str("value", string(m.Value)).
		Msg("consume")

	return nil
}
