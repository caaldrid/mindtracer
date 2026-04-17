package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/caaldrid/mindtracer/backend/auth"
	"github.com/caaldrid/mindtracer/backend/handlers"
	"github.com/caaldrid/mindtracer/backend/setup"
	"github.com/caaldrid/mindtracer/backend/storage"
)

func corsMiddleware(allowedOrigin string) gin.HandlerFunc {
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

// New creates a gin.Engine with all middleware, route groups, and handlers wired up.
func New(c setup.Config, store storage.Storage) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware("*"))

	authGroup := router.Group("/api/auth")
	account := handlers.NewAccountHandler(store.Users, c.SecretKey, c.TokenLifespan)
	account.RegisterRoutes(authGroup)

	authorized := router.Group("/api/")
	authorized.Use(auth.JWTMiddleware(c.SecretKey))

	return router
}
