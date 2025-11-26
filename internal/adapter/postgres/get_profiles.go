package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.noway/internal/domain"
	"gitlab.noway/internal/dto"
	"gitlab.noway/pkg/otel/tracer"
	"gitlab.noway/pkg/transaction"
)

type GetProfilesDTO struct {
	ID        pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	DeletedAt pgtype.Timestamptz
	Name      pgtype.Text
	Age       pgtype.Int4
	Status    pgtype.Text
	Verified  pgtype.Bool
	Contacts  []byte
}

func (d *GetProfilesDTO) ToDomain() (domain.Profile, error) {
	var contacts domain.Contacts

	err := json.Unmarshal(d.Contacts, &contacts)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("cannot unmarshal contacts: %w", err)
	}

	return domain.Profile{
		ID:        d.ID.Bytes,
		CreatedAt: d.CreatedAt.Time,
		UpdatedAt: d.UpdatedAt.Time,
		DeletedAt: d.DeletedAt.Time,
		Name:      domain.Name(d.Name.String),
		Age:       domain.Age(d.Age.Int32),
		Status:    domain.NewStatus(d.Status.String),
		Verified:  d.Verified.Bool,
		Contacts:  contacts,
	}, nil
}

func (d *GetProfilesDTO) Dest() []any {
	return []any{
		&d.ID,
		&d.CreatedAt,
		&d.UpdatedAt,
		&d.DeletedAt,
		&d.Name,
		&d.Age,
		&d.Status,
		&d.Verified,
		&d.Contacts,
	}
}

func (p *Postgres) GetProfiles(ctx context.Context, input dto.GetProfilesInput) ([]domain.Profile, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres GetProfiles")
	defer span.End()

	sql := `SELECT id, created_at, updated_at, deleted_at, name, age, status, verified, contacts
                 FROM profile
                 ORDER BY %s %s
                 OFFSET %d
                 LIMIT %d`

	sql = fmt.Sprintf(sql, input.Sort, input.Order, input.Offset, input.Limit)

	txOrPool := transaction.TryExtractTX(ctx)

	rows, err := txOrPool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("txOrPool.QueryRow: %w", err)
	}

	defer rows.Close()

	profiles := make([]domain.Profile, 0, input.Limit)

	for rows.Next() {
		var d GetProfilesDTO

		err = rows.Scan(d.Dest()...)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		profile, err := d.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("dto.ToDomain: %w", err)
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
