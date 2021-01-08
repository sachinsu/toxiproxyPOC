package main

import (
	"fmt"

	"github.com/sachinsu/chaosengg/api/server"
	"github.com/sachinsu/chaosengg/internal/setup"
)

func main() {
	dbConnString := "postgres://postgres:@localhost:32777/proxytest?sslmode=disable"

	tproxy, err := setup.AddPGProxy("[::]:32777", "localhost:5432")
	if err != nil {
		fmt.Printf("Error with toxiproxy %v\n", err)
		return
	}

	defer tproxy.Delete()

	// perform db migrations
	err = setup.MigrateDBSchema(dbConnString, "file://assets")

	if err != nil {
		fmt.Printf("Error with database setup %v\n", err)
		return
	}
	// run HTTP server
	Srvr, err := server.NewServer(dbConnString, ":8080", "/api")

	if err != nil {
		fmt.Printf("Error with HTTP Server %v\n", err)
	}

	err = Srvr.Run()

	if err != nil {
		fmt.Printf("Error with HTTP Server %v\n", err)
	}
}
