package postgres

import (
	"context"
	"fmt"

	"gitlab.noway/internal/domain"
	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/transaction"
)

func (p *Postgres) CreateProperty(ctx context.Context, property domain.Property) error {
	ctx, span := tracer.Start(ctx, "adapter postgres CreateProperty")
	defer span.End()

	const sql = `INSERT INTO property (profile_id, tags)
                    VALUES ($1, $2)`

	args := []any{
		property.ProfileID,
		property.Tags,
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err := txOrPool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
