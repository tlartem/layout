package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/segmentio/kafka-go"

	"gitlab.noway/internal/domain"
	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/transaction"
)

func (p *Postgres) SaveOutboxKafka(ctx context.Context, msgs ...kafka.Message) error {
	ctx, span := tracer.Start(ctx, "adapter postgres SaveOutboxKafka")
	defer span.End()

	if len(msgs) == 0 {
		return nil
	}

	batch := make([]any, 0, len(msgs))

	for _, msg := range msgs {
		if msg.Topic == "" {
			return domain.ErrEmptyTopic
		}

		headers, err := json.Marshal(msg.Headers)
		if err != nil {
			return fmt.Errorf("json.Marshal headers: %w", err)
		}

		batch = append(batch, goqu.Record{
			"topic":   msg.Topic,
			"key":     msg.Key,
			"value":   msg.Value,
			"headers": headers,
		})
	}

	sql, _, err := goqu.Insert("outbox").Rows(batch...).ToSQL()
	if err != nil {
		return fmt.Errorf("goqu.Insert.ToSQL: %w", err)
	}

	txOrPool := transaction.TryExtractTX(ctx)

	_, err = txOrPool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}
