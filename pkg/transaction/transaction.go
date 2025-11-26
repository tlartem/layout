package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/postgres"
)

var (
	errMissingInit  = errors.New("missing `transaction.Init' call before `transaction.Begin'")
	errMissingBegin = errors.New("missing `transaction.Begin' call before 'transaction.Get'")
)

//nolint:gochecknoglobals
var (
	pool       *pgxpool.Pool
	IsUnitTest bool
)

type ctxKey struct{}

func Init(p *postgres.Pool) {
	pool = p.Pool
}

type Transaction struct {
	pgx.Tx
}

func Begin(ctx context.Context) (context.Context, error) {
	ctx, span := tracer.Start(ctx, "transaction Begin")
	defer span.End()

	if IsUnitTest {
		return ctx, nil
	}

	if pool == nil {
		return nil, errMissingInit
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Begin: %w", err)
	}

	ctx = context.WithValue(ctx, ctxKey{}, &Transaction{tx})

	return ctx, nil
}

func Rollback(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "transaction Rollback")
	defer span.End()

	if IsUnitTest {
		return
	}

	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return
	}

	err := tx.Rollback(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		log.Error().Err(err).Msg("transaction: Rollback")
	}
}

func Commit(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "transaction Commit")
	defer span.End()

	if IsUnitTest {
		return nil
	}

	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return errMissingBegin
	}

	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func TryExtractTX(ctx context.Context) Executor { //nolint:ireturn
	tx, ok := ctx.Value(ctxKey{}).(*Transaction)
	if !ok {
		return pool
	}

	return tx
}
