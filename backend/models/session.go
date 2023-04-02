package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID string `json:"id"`
	// Foregin key to user
	UserID  string    `gorm:"type:uuid;not null" json:"user_id"`
	Expires time.Time `json:"expires"`
}

func (s *Session) IsExpired() bool {
	return s.Expires.Before(time.Now())
}

func (s *Session) BeforeCreate() error {
	s.ID = uuid.NewString()
	return nil
}
