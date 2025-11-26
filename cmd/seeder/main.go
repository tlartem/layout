package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/profile"
)

const (
	iterates  = 1_000
	batchSize = 1_000
	dbURL     = "postgres://login:pass@localhost:5432/postgres"
)

// go tool pprof -http=:8080 -base=a_cpu.pprof b_cpu.pprof

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	now := time.Now()

	ctx := context.Background()
	Seeder(ctx, dbURL)

	fmt.Println("Done:", time.Since(now))
}

func Seeder(ctx context.Context, dbURL string) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("pgx.Connect: %v\n", err)
	}

	for range iterates {
		ids := GenIDs()

		profileBatch := ProfileBatch(ids)

		err = Insert(ctx, conn, profileBatch, "profile")
		if err != nil {
			log.Fatalf("profile: Insert: %v\n", err)
		}

		propertyBatch := PropertyBatch(ids)

		err = Insert(ctx, conn, propertyBatch, "property")
		if err != nil {
			log.Fatalf("property: Insert: %v\n", err)
		}
	}
}

func ProfileBatch(ids []uuid.UUID) []any {
	batch := make([]any, 0, batchSize)

	for _, id := range ids {
		batch = append(batch, goqu.Record{
			"id":       id,
			"name":     gofakeit.Name(),
			"age":      gofakeit.IntRange(18, 120),
			"status":   Status(),
			"verified": gofakeit.Bool(),
			"contacts": Contacts(),
		})
	}

	return batch
}

func PropertyBatch(ids []uuid.UUID) []any {
	batch := make([]any, 0, batchSize)

	for _, id := range ids {
		batch = append(batch, goqu.Record{
			"profile_id": id,
			"tags":       Tags(),
		})
	}

	return batch
}
