package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
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
