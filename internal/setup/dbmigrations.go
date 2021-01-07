package setup

import (
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
		log.Fatal(err)
		return err
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
		return err
	}
}
