package main

import (
	"log"

	"github.com/caaldrid/mindtracer/backend/api"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func main() {
	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}
	api.StartServer(config)
}
