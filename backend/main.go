package main

import (
	"flag"
	"log"

	"github.com/caaldrid/mindtracer/backend/api"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func main() {
	runMigration := flag.Bool(
		"migrate",
		false,
		"Tells us to run migration instead of starting the server",
	)

	flag.Parse()

	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	DB, err := setup.ConnectDB(setup.FormatDBConnectionString(config))
	if err != nil {
		log.Fatal("? Could connect to database instance", err)
	}

	if *runMigration {
		setup.MigrateModels(DB)
	} else {
		api.StartServer(config, DB)
	}
}
