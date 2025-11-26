package repository

import (
	"time"

	"gitlab.noway/internal/adapter/postgres"
	"gitlab.noway/pkg/redis"
)

const (
	prefix = "noway:layout:"
	ttl    = time.Minute
)

type Repository struct {
	redis    *redis.Client
	postgres *postgres.Postgres
}

func New(client *redis.Client, p *postgres.Postgres) *Repository {
	return &Repository{
		redis:    client,
		postgres: p,
	}
}
