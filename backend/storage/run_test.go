package storage_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/setup"
)

var TestDB *gorm.DB

func randomString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	postgresContainer, _ := postgres.Run(ctx,
		"postgres:18.3-trixie",
		postgres.WithDatabase("mindtracer"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword(randomString(16)),
		postgres.BasicWaitStrategies(),
	)

	cString, _ := postgresContainer.ConnectionString(ctx,
		"sslmode=disable",
		"TimeZone=UTC",
	)

	DB, err := setup.ConnectDB(cString)
	if err != nil {
		log.Fatal("Could not connect to database instance", err)
	}

	setup.MigrateModels(DB)

	TestDB = DB

	code := m.Run()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}
