package setup

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateDBSchema - applies migrations to DB
func MigrateDBSchema(dbConn string, migrationsPath string) error {
	// m, err := migrate.New(
	// 	"file://assets",
	// 	"postgres://postgres:postgres@localhost:5432/example?sslmode=disable")

	m, err := migrate.New(
		migrationsPath,
		dbConn)

	if err != nil {
		log.Fatalf("Error while initiating db Migration %v", err)
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error during applying migrations %v", err)
		return err
	}

	return nil
}
