package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"

	"gitlab.noway/pkg/otel/tracer"
)

func Wrap(ctx context.Context, fn func(context.Context) error) error {
	ctx, span := tracer.Start(ctx, "transaction Wrap")
	defer span.End()

	if IsUnitTest {
		return fn(ctx)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin: %w", err)
	}

	defer func() {
		err = tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Msg("transaction: Rollback")
		}
	}()

	ctx = context.WithValue(ctx, ctxKey{}, &Transaction{tx})

	err = fn(ctx)
	if err != nil {
		return fmt.Errorf("fn: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}
