package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrAreaNotFound      = errors.New("area not found")
	ErrAreaAlreadyExists = errors.New("area with that name already exists for this user")
)

type AreaStorage interface {
	Create(ctx context.Context, area *models.Area) error
	FindByID(ctx context.Context, userID, areaID uuid.UUID) (*models.Area, error)
	ListByUser(ctx context.Context, userID uuid.UUID, archived bool) ([]models.Area, error)
	Update(ctx context.Context, area *models.Area) error
	Archive(ctx context.Context, userID, areaID uuid.UUID) error
	Delete(ctx context.Context, userID, areaID uuid.UUID) error
}

type areaStorage struct {
	db *gorm.DB
}

func NewAreaStorage(db *gorm.DB) AreaStorage {
	return &areaStorage{db: db}
}

func (s *areaStorage) Create(ctx context.Context, area *models.Area) error {
	if err := s.db.WithContext(ctx).Create(area).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAreaAlreadyExists
		}
		return fmt.Errorf("areaStorage.Create: %w", err)
	}
	return nil
}

func (s *areaStorage) FindByID(ctx context.Context, userID, areaID uuid.UUID) (*models.Area, error) {
	var area models.Area
	err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", areaID, userID).First(&area).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAreaNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("areaStorage.FindByID: %w", err)
	}
	return &area, nil
}

func (s *areaStorage) ListByUser(ctx context.Context, userID uuid.UUID, archived bool) ([]models.Area, error) {
	var areas []models.Area
	if err := s.db.WithContext(ctx).Where("user_id = ? AND is_archived = ?", userID, archived).Find(&areas).Error; err != nil {
		return nil, fmt.Errorf("areaStorage.ListByUser: %w", err)
	}
	return areas, nil
}

func (s *areaStorage) Update(ctx context.Context, area *models.Area) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", area.ID, area.UserID).Save(area)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return ErrAreaAlreadyExists
		}
		return fmt.Errorf("areaStorage.Update: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAreaNotFound
	}
	return nil
}

func (s *areaStorage) Archive(ctx context.Context, userID, areaID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Model(&models.Area{}).
		Where("id = ? AND user_id = ?", areaID, userID).
		Update("is_archived", true)
	if result.Error != nil {
		return fmt.Errorf("areaStorage.Archive: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAreaNotFound
	}
	return nil
}

func (s *areaStorage) Delete(ctx context.Context, userID, areaID uuid.UUID) error {
	result := s.db.WithContext(ctx).Unscoped().Where("id = ? AND user_id = ?", areaID, userID).Delete(&models.Area{})
	if result.Error != nil {
		return fmt.Errorf("areaStorage.Delete: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAreaNotFound
	}
	return nil
}
