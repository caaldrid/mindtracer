package handlers_test

import (
	"context"

	"github.com/caaldrid/mindtracer/backend/models"
)

type mockUserStorage struct {
	createIfNotExistsErr error
	findByUsernameUser   *models.User
	findByUsernameErr    error
}

func (m *mockUserStorage) CreateIfNotExists(_ context.Context, _ *models.User) error {
	return m.createIfNotExistsErr
}

func (m *mockUserStorage) FindByUsername(_ context.Context, _ string) (*models.User, error) {
	return m.findByUsernameUser, m.findByUsernameErr
}
