package entity

import (
	"github.com/google/uuid"
	"time"
)

type Model struct {
	UUID      uuid.UUID `gorm:"type:uuid;primaryKey;default:(gen_random_uuid())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type ClassID uint8

type Roadmap uint8

type Row struct {
	Model
	ClassID ClassID `gorm:"index:idx_class_id"`
	Roadmap Roadmap
}
