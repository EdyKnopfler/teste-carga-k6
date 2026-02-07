package amigos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
}

type Amigo struct {
	Base
	Nome           string
	DataNascimento time.Time
	Preferencias   []Preferencia `gorm:"foreignKey:IDAmigo"` // 1:N
}

type Preferencia struct {
	Base
	IDAmigo uuid.UUID `gorm:"type:uuid;"`
	Nome    string
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	u6, err := uuid.NewV6()

	if err != nil {
		return err
	}

	b.ID = u6
	return nil
}
