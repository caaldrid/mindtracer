package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceType string

const (
	ResourceTypeArticle ResourceType = "Article"
	ResourceTypeBook    ResourceType = "Book"
	ResourceTypeVideo   ResourceType = "Video"
	ResourceTypeNote    ResourceType = "Note"
)

type Resource struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE"`
	Title       string
	Description string
	ISBN        *string `gorm:"check:only_books,ISBN is NULL OR type = 'Book'"`
	SourceURL   *string
	Type        ResourceType
	IsArchived  bool
}

type ProjectResource struct {
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;primary_key"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null;constraint:OnDelete:CASCADE;primary_key"`
	LinkedAt   time.Time
}
