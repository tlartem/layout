//nolint:wrapcheck
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
	"gitlab.noway/pkg/otel/tracer"
)

func (r *Repository) ReadOutboxKafka(ctx context.Context, limit int) ([]kafka.Message, error) {
	ctx, span := tracer.Start(ctx, "adapter repository ReadOutboxKafka")
	defer span.End()

	return r.postgres.ReadOutboxKafka(ctx, limit)
}

func (r *Repository) SaveOutboxKafka(ctx context.Context, msgs ...kafka.Message) error {
	ctx, span := tracer.Start(ctx, "adapter repository SaveOutboxKafka")
	defer span.End()

	return r.postgres.SaveOutboxKafka(ctx, msgs...)
}

func (r *Repository) CreateProfile(ctx context.Context, profile domain.Profile) error {
	ctx, span := tracer.Start(ctx, "adapter repository CreateProfile")
	defer span.End()

	return r.postgres.CreateProfile(ctx, profile)
}

func (r *Repository) CreateProperty(ctx context.Context, property domain.Property) error {
	ctx, span := tracer.Start(ctx, "adapter repository CreateProperty")
	defer span.End()

	return r.postgres.CreateProperty(ctx, property)
}

func (r *Repository) GetProfile(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	ctx, span := tracer.Start(ctx, "adapter repository GetProfile")
	defer span.End()

	profile, err := r.getCache(ctx, profileID)
	if err != nil { //nolint:nestif
		if errors.Is(err, domain.ErrNotFound) {
			profile, err = r.postgres.GetProfile(ctx, profileID)
			if err != nil {
				return domain.Profile{}, fmt.Errorf("r.postgres.GetProfile: %w", err)
			}

			// Закомментить
			if profile.IsDeleted() {
				return profile, domain.ErrNotFound
			}

			err = r.setCache(ctx, profile)
			if err != nil {
				log.Error().Err(err).Str("profileID", profileID.String()).Msg("cache: GetProfile: set cache")
			}

			return profile, nil
		}

		log.Error().Err(err).Str("profileID", profileID.String()).Msg("cache: GetProfile: get cache")

		return r.postgres.GetProfile(ctx, profileID)
	}

	return profile, nil
}

func (r *Repository) GetProfiles(ctx context.Context, input dto.GetProfilesInput) ([]domain.Profile, error) {
	ctx, span := tracer.Start(ctx, "adapter repository GetProfiles")
	defer span.End()

	return r.postgres.GetProfiles(ctx, input)
}

func (r *Repository) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	ctx, span := tracer.Start(ctx, "adapter repository UpdateProfile")
	defer span.End()

	err := r.postgres.UpdateProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("r.postgres.UpdateProfile: %w", err)
	}

	err = r.setCache(ctx, profile)
	if err != nil {
		log.Error().Err(err).Str("profileID", profile.ID.String()).Msg("cache: UpdateProfile: update cache")
	}

	return nil
}

func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "adapter repository DeleteProfile")
	defer span.End()

	err := r.postgres.DeleteProfile(ctx, id)
	if err != nil {
		return fmt.Errorf("r.postgres.DeleteProfile: %w", err)
	}

	err = r.deleteCache(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("profileID", id.String()).Msg("cache: DeleteProfile: delete cache")
	}

	return nil
}
