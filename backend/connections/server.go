package connections

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/caaldrid/mindtracer/backend/setup"
)

func StartServer(c setup.Config) {
	db, err := setup.ConnectDB(c)
	if err != nil {
		log.Fatal("? Could connect to database instance", err)
	}

	router := gin.Default()

	setupAccountHandler(db, router)

	if err := router.Run(fmt.Sprintf(":%s", c.ServerPort)); err != nil {
		log.Fatal("Failed to start router", err)
	}
}
