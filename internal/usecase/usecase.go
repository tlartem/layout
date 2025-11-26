package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
	"gitlab.noway/pkg/httpclient"
)

//go:generate mockery

type Redis interface {
	IsExists(ctx context.Context, idempotencyKey string) bool
}

type Kafka interface {
	Produce(ctx context.Context, msgs ...kafka.Message) error
}

type Profile interface {
	Create(ctx context.Context, name string, age int, email, phone string) (uuid.UUID, error)
	Delete(ctx context.Context, id string) error
	GetProfile(ctx context.Context, id string) (httpclient.Profile, error)
	Update(ctx context.Context, id string, name *string, age *int, email, phone *string) error
}

type Postgres interface {
	CreateProfile(ctx context.Context, profile domain.Profile) error
	CreateProperty(ctx context.Context, property domain.Property) error
	GetProfile(ctx context.Context, profileID uuid.UUID) (domain.Profile, error)
	GetProfiles(ctx context.Context, input dto.GetProfilesInput) ([]domain.Profile, error)
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error

	ReadOutboxKafka(ctx context.Context, limit int) ([]kafka.Message, error)
	SaveOutboxKafka(ctx context.Context, msgs ...kafka.Message) error
}

type UseCase struct {
	profile  Profile
	postgres Postgres
	kafka    Kafka
	redis    Redis
}

func New(postgres Postgres, profile Profile, k Kafka, redis Redis) *UseCase {
	return &UseCase{
		postgres: postgres,
		profile:  profile,
		kafka:    k,
		redis:    redis,
	}
}
