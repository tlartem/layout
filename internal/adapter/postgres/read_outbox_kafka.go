package postgres

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/transaction"
)

func (p *Postgres) ReadOutboxKafka(ctx context.Context, limit int) ([]kafka.Message, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres ReadOutboxKafka")
	defer span.End()

	const sql = `WITH taken AS (SELECT id, topic, key, value, headers
					   FROM outbox
					   ORDER BY created_at
					   LIMIT $1 FOR UPDATE SKIP LOCKED)
				DELETE
				FROM outbox
				WHERE id IN (SELECT id FROM taken)
				RETURNING topic, key, value, headers;`

	txOrPool := transaction.TryExtractTX(ctx)

	rows, err := txOrPool.Query(ctx, sql, limit)
	if err != nil {
		return nil, fmt.Errorf("txOrPool.Query: %w", err)
	}

	defer rows.Close()

	msgs := make([]kafka.Message, 0, limit)

	for rows.Next() {
		var msg kafka.Message

		err = rows.Scan(&msg.Topic, &msg.Key, &msg.Value, &msg.Headers)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}
