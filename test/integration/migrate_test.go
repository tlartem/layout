//go:build integration

package test

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (s *Suite) ResetMigrations() {
	const (
		filesPath = "file://../../migration/postgres"
		dbURL     = "postgres://login:pass@localhost:5432/postgres?sslmode=disable"
	)

	m, err := migrate.New(filesPath, dbURL)
	s.NoError(err)

	err = m.Down()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	}
	s.NoError(err)

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	}
	s.NoError(err)
}
