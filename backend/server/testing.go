package server

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"

	"github.com/caaldrid/mindtracer/backend/setup"
	"github.com/caaldrid/mindtracer/backend/storage"
)

// NewTestServer creates a gin.Engine in test mode with full routing, middleware,
// and an internally generated config. Only storage needs to be provided.
func NewTestServer(store storage.Storage) *gin.Engine {
	gin.SetMode(gin.TestMode)

	b := make([]byte, 32)
	rand.Read(b)
	secret := hex.EncodeToString(b)

	config := setup.Config{
		SecretKey:     secret,
		TokenLifespan: 4,
	}

	return New(config, store)
}
