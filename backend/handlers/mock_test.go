package handlers_test

import (
	"context"

	"github.com/caaldrid/mindtracer/backend/models"
)

type mockUserStorage struct {
	createErr          error
	findByUsernameUser *models.User
	findByUsernameErr  error
}

func (m *mockUserStorage) Create(_ context.Context, _ *models.User) error {
	return m.createErr
}

func (m *mockUserStorage) FindByUsername(_ context.Context, _ string) (*models.User, error) {
	return m.findByUsernameUser, m.findByUsernameErr
}
