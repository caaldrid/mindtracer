package api_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/caaldrid/mindtracer/backend/api"
	"github.com/caaldrid/mindtracer/backend/setup"
)

func randomString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var TestServer *httptest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()
	c := setup.Config{
		DBUserName:     "postgres",
		DBUserPassword: randomString(16),
		DBName:         "mindtracer",
		SecretKey:      randomString(32),
		TokenLifespan:  4,
	}

	postgresContainer, _ := postgres.Run(ctx,
		"postgres:18.3-trixie",
		postgres.WithDatabase(c.DBName),
		postgres.WithUsername(c.DBUserName),
		postgres.WithPassword(c.DBUserPassword),
		postgres.BasicWaitStrategies(),
	)

	cString, _ := postgresContainer.ConnectionString(ctx,
		"sslmode=disable",
		"TimeZone=UTC",
	)

	DB, err := setup.ConnectDB(cString)
	if err != nil {
		log.Fatal("? Could connect to database instance", err)
	}

	setup.MigrateModels(DB)
	router := api.ConfigRouter(c, DB)

	TestServer = httptest.NewServer(router)

	code := m.Run()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	TestServer.Close()

	os.Exit(code)
}
