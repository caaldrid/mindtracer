package storage

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
)

var ErrUserAlreadyExists = errors.New("user already exists for the given username")

type UserStorage interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	CreateIfNotExists(ctx context.Context, user *models.User) error
}

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStorage) CreateIfNotExists(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.User
		err := tx.Where("user_name = ?", user.UserName).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			return ErrUserAlreadyExists
		}
		return tx.Create(user).Error
	})
}
