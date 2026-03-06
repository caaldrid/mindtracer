package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	AreaID         uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	PrerequisiteID *uuid.UUID `gorm:"type:uuid;constraint:OnDelete:SET NULL"`
	Name           string
	Description    string
	IsArchived     bool
}
