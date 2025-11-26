package redis

import (
	"time"

	"gitlab.noway/pkg/redis"
)

const (
	idempotencyPrefix = "noway:layout:idempotency:"
	ttl               = time.Hour
)

type Redis struct {
	redis *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{
		redis: client,
	}
}
