package handlers_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/server"
	"github.com/caaldrid/mindtracer/backend/storage"
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

func randomString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func TestRegister(t *testing.T) {
	username := randomString(8)
	email := fmt.Sprintf("%s@test.com", randomString(8))
	password := randomString(16)

	tests := []struct {
		name       string
		body       string
		mock       *mockUserStorage
		wantStatus int
		wantBody   string
	}{
		{
			name:       "successful registration",
			body:       fmt.Sprintf(`{"username":"%s","password":"%s","email":"%s"}`, username, password, email),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusCreated,
			wantBody:   "has been registered",
		},
		{
			name:       "missing username",
			body:       fmt.Sprintf(`{"password":"%s","email":"%s"}`, password, email),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       fmt.Sprintf(`{"username":"%s","email":"%s"}`, username, email),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       `{}`,
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate email",
			body:       fmt.Sprintf(`{"username":"%s","password":"%s","email":"%s"}`, username, password, email),
			mock:       &mockUserStorage{createErr: storage.ErrUserAlreadyExists},
			wantStatus: http.StatusConflict,
			wantBody:   "user with that email already exists",
		},
		{
			name:       "storage error",
			body:       fmt.Sprintf(`{"username":"%s","password":"%s","email":"%s"}`, username, password, email),
			mock:       &mockUserStorage{createErr: errors.New("db connection lost")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.Storage{Users: tt.mock}
			router := server.NewTestServer(store)

			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantBody != "" && !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("body %q does not contain %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	password := randomString(16)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	validUser := &models.User{
		ID:       uuid.New(),
		UserName: randomString(8),
		Email:    fmt.Sprintf("%s@test.com", randomString(8)),
		Password: string(hashedPassword),
	}

	tests := []struct {
		name       string
		body       string
		mock       *mockUserStorage
		wantStatus int
		wantBody   string
	}{
		{
			name:       "successful login",
			body:       fmt.Sprintf(`{"email":"%s","password":"%s"}`, validUser.Email, password),
			mock:       &mockUserStorage{findByEmailUser: validUser},
			wantStatus: http.StatusOK,
			wantBody:   "accessToken",
		},
		{
			name:       "missing email",
			body:       fmt.Sprintf(`{"password":"%s"}`, password),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing password",
			body:       fmt.Sprintf(`{"email":"%s"}`, validUser.Email),
			mock:       &mockUserStorage{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "email not found",
			body:       fmt.Sprintf(`{"email":"unknown@test.com","password":"%s"}`, password),
			mock:       &mockUserStorage{findByEmailErr: errors.New("not found")},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Invalid email",
		},
		{
			name:       "wrong password",
			body:       fmt.Sprintf(`{"email":"%s","password":"%s"}`, validUser.Email, randomString(16)),
			mock:       &mockUserStorage{findByEmailUser: validUser},
			wantStatus: http.StatusUnauthorized,
			wantBody:   "Invalid Password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.Storage{Users: tt.mock}
			router := server.NewTestServer(store)

			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if tt.wantBody != "" && !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("body %q does not contain %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}
