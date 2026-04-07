package setup

import (
	"log"

	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
)

func MigrateModels(DB *gorm.DB) {
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Area{})
	DB.AutoMigrate(&models.Resource{})
	DB.AutoMigrate(&models.Project{})
	DB.AutoMigrate(&models.ToDo{})
	DB.AutoMigrate(&models.ProjectResource{})

	log.Println("Migration complete")
}
