package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/log15adapter"
	_ "github.com/jackc/pgx/v4"
)

type City struct {
	ID      int
	Name    string
	Website string
	Mayor   string
}

func openDb,(ctx context.Context, dbConn string) (*pgx.Conn, func(context.Context) error, error) {
	config, err := pgx.ParseConfig(dbConn)
	if err != nil {
		return nil, nil, fmt.Errorf("Error with database connection string - %w", err)
	}

	config.Logger = log15adapter.NewLogger(log.New("dbfetcher", "pgx"))

	db, err := pgx.Connect(ctx, config)

	if err != nil {
		log.Fatal(err)
		return nil, nil, fmt.Errorf("Error opening db connection - %w", err)
	}

	return db, db.Close, nil
}

// GetCities - from DB
func GetCities(ctx context.Context, dbConn string, start int, count int) ([]City, error) {
	var cities []City

	db, close, err := openDb(ctx, dbConn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer close(ctx)

	rows, err := db.Query("SELECT id, cityname, website, mayor FROM cities LIMIT $1", count)

	if err != nil {
		log.Println("Error while retrieving rows %v", err)
		return nil, fmt.Errorf("Error while retrieving rows %w", err)
	}

	for rows.Next() {
		newCity := City{}

		if err := rows.Scan(&City.Id, &City.Name, &City.Website, &City.Mayor); err != nil {
			return nil, err
		}
		cities = append(cities, newCity)
	}

	return cities, rows.Err()

}
