package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	StatusInactive Status = "Inactive"
	StatusWorking  Status = "Working"
	StatusBlocked  Status = "Blocked"
	StatusClosed   Status = "Closed"
)

type ToDo struct {
	gorm.Model
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ProjectID      uuid.UUID  `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	PrerequisiteID *uuid.UUID `gorm:"type:uuid;constraint:OnDelete:SET NULL"`
	Status         Status     `gorm:"default:'Inactive'"`
	Title          string
	Description    string
	DueDate        *time.Time
	CompletedAt    *time.Time
}
