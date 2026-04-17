package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/caaldrid/mindtracer/backend/setup"
	"github.com/caaldrid/mindtracer/backend/storage"
)

func CORSMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func ConfigRouter(c setup.Config, store storage.Storage) *gin.Engine {
	router := gin.Default()
	router.Use(CORSMiddleware("*"))

	setupAccountHandler(store.Users, router, c)

	authorized := router.Group("/api/")
	authorized.Use(jwtAuthMiddleware(c.SecretKey))

	return router
}

func StartServer(c setup.Config, store storage.Storage) {
	router := ConfigRouter(c, store)

	if err := router.Run(fmt.Sprintf(":%s", c.ServerPort)); err != nil {
		log.Fatal("Failed to start router", err)
	}
}
