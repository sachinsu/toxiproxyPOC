package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

// City structure
type City struct {
	ID      int
	Name    string
	Website string
	Mayor   string
}

func openDb(ctx context.Context, dbConn string) (*pgx.Conn, func(context.Context) error, error) {
	config, err := pgx.ParseConfig(dbConn)
	if err != nil {
		return nil, nil, fmt.Errorf("Error with database connection string - %w", err)
	}

	// config.Logger = log15adapter.NewLogger(log.New(os.Stdout, "database: ", log.Ldate|log.Ltime|log.Lshortfile))

	db, err := pgx.ConnectConfig(ctx, config)

	if err != nil {
		log.Fatal(err)
		return nil, nil, fmt.Errorf("Error opening db connection - %w", err)
	}

	return db, db.Close, nil
}

// GetCities - from DB
func GetCities(ctx context.Context, dbConn string, startId int, count int) ([]City, error) {
	var cities []City

	db, close, err := openDb(ctx, dbConn)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer close(ctx)

	var rows pgx.Rows

	if startId > 0 {
		rows, err = db.Query(ctx, `SELECT id, cityname, website, mayor FROM cities 
				where id > $1 
				order by id
				LIMIT $2`, startId, count)
	} else {
		rows, err = db.Query(ctx, `SELECT id, cityname, website, mayor FROM cities 
				 order by id LIMIT $1`, count)
	}

	if err != nil {
		log.Printf("Error while retrieving rows %v\n", err)
		return nil, fmt.Errorf("Error while retrieving rows %w", err)
	}

	for rows.Next() {
		newCity := City{}

		if err := rows.Scan(&newCity.ID, &newCity.Name, &newCity.Website, &newCity.Mayor); err != nil {
			return nil, err
		}
		cities = append(cities, newCity)
	}

	return cities, rows.Err()

}
