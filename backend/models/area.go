package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Name        string
	Description string
	IsArchived  bool
}
