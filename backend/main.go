package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/caaldrid/mindtracer/backend/server"
	"github.com/caaldrid/mindtracer/backend/setup"
	"github.com/caaldrid/mindtracer/backend/storage"
)

func main() {
	runMigration := flag.Bool(
		"migrate",
		false,
		"Runs DB migration instead of starting the server.",
	)

	seedDatabase := flag.Bool("seed", false, "Seed the database with fixture data.")
	teardownDatabase := flag.Bool("teardown", false, "Truncate all tables and exit.")

	flag.Parse()

	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	DB, err := setup.ConnectDB(setup.FormatDBConnectionString(config))
	if err != nil {
		log.Fatal("? Could connect to database instance", err)
	}

	cntx := context.Background()

	switch {
	case *runMigration:
		setup.MigrateModels(DB)
	case *seedDatabase:
		if err := setup.SeedDB(cntx, DB); err != nil {
			log.Fatalf("Seed failed: %s", err)
		}
	case *teardownDatabase:
		if err := setup.TeardownDB(cntx, DB); err != nil {
			log.Fatalf("Teardown failed: %s", err)
		}
	default:
		store := storage.Storage{
			Users: storage.NewUserStorage(DB),
			Areas: storage.NewAreaStorage(DB),
		}
		router := server.New(config, store)
		if err := router.Run(fmt.Sprintf(":%s", config.ServerPort)); err != nil {
			log.Fatal("Failed to start server", err)
		}
	}
}
