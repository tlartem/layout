package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.noway/internal/domain"
	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/transaction"
)

func (p *Postgres) UpdateProfile(ctx context.Context, profile domain.Profile) error {
	ctx, span := tracer.Start(ctx, "adapter postgres UpdateProfile")
	defer span.End()

	const sql = `UPDATE profile SET name = $1, age = $2, contacts = $3, updated_at = NOW()
                     WHERE id = $4`

	contacts, err := json.Marshal(profile.Contacts)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	args := []any{
		profile.Name,
		profile.Age,
		contacts,
		profile.ID,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err = txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}

		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
