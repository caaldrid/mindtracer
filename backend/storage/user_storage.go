package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrUserAlreadyExists = errors.New("user with that email already exists")

type UserStorage interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStorage) Create(ctx context.Context, user *models.User) error {
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("userStorage.Create: %w", err)
	}
	return nil
}
