package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
}

// BeforeCreate hook to set the UUID before creating a new record
func (b *BaseModel) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New() // Generates a UUIDv4
	}
	return nil
}
