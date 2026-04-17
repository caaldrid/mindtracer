package handlers_test

import (
	"context"

	"github.com/caaldrid/mindtracer/backend/models"
)

type mockUserStorage struct {
	createErr       error
	findByEmailUser *models.User
	findByEmailErr  error
}

func (m *mockUserStorage) Create(_ context.Context, _ *models.User) error {
	return m.createErr
}

func (m *mockUserStorage) FindByEmail(_ context.Context, _ string) (*models.User, error) {
	return m.findByEmailUser, m.findByEmailErr
}
