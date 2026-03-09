package main

import (
	"log"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func main() {
	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	DB, err := setup.ConnectDB(config)
	if err != nil {
		log.Fatal("? Could connect to database instance", err)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Area{})
	DB.AutoMigrate(&models.Resource{})
	DB.AutoMigrate(&models.Project{})
	DB.AutoMigrate(&models.ToDo{})
	DB.AutoMigrate(&models.ProjectResource{})

	log.Println("? Migration complete")
}
