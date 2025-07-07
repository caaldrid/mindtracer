package main

import (
	"log"

	"github.com/caaldrid/mindtracer/backend/connections"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func main() {
	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	connections.StartServer(config)
}
