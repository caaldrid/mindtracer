package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/setup"
)

func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Header(
			"Access-Control-Allow-Origin",
			allowedOrigin,
		) // Change '*' to specific origins in production
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent) // No content response for OPTIONS
			return
		}

		c.Next() // Proceed to the next handler
	}
}

func ConfigRouter(c setup.Config, DB *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware("*"))

	setupAccountHandler(DB, router, c)

	authorized := router.Group("/api/")
	authorized.Use(jwtAuthMiddleware(c.SecretKey))

	return router
}

func StartServer(c setup.Config, DB *gorm.DB) {
	router := ConfigRouter(c, DB)

	if err := router.Run(fmt.Sprintf(":%s", c.ServerPort)); err != nil {
		log.Fatal("Failed to start router", err)
	}
}
