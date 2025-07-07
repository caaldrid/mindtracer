package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserName string    `gorm:"type:varchar(255);not null"`
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
}
