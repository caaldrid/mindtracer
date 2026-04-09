package handlers_test

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/caaldrid/mindtracer/backend/setup"
)

func randomString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var testConfig = setup.Config{
	SecretKey:     randomString(32),
	TokenLifespan: 4,
}
