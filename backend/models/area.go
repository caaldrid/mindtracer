package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;uniqueIndex:idx_user_area_name"`
	Name        string    `gorm:"uniqueIndex:idx_user_area_name"`
	Description string
	IsArchived  bool
}
