package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"gitlab.noway/pkg/otel/tracer"
)

func (r *Redis) IsExists(ctx context.Context, idempotencyKey string) bool {
	ctx, span := tracer.Start(ctx, "adapter redis IsExists")
	defer span.End()

	key := idempotencyPrefix + idempotencyKey

	err := r.redis.Get(ctx, key).Err()
	if err == nil { // err == nil
		return true
	}

	if !errors.Is(err, redis.Nil) {
		log.Error().Err(err).Msg("redis: IsExists: r.redis.Get")
	}

	err = r.redis.Set(ctx, key, []byte{}, ttl).Err()
	if err != nil {
		log.Error().Err(err).Msg("redis: IsExists: r.redis.Set")
	}

	return false
}
