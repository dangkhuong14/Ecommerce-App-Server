package domain

import (
	"ecommerce/common"
	"time"
)

type Session struct {
	id           common.UUID
	userID       common.UUID
	refreshToken string
	accessExpAt  time.Time
	refreshExpAt time.Time
}

// Constructor
func NewSession(id, userID common.UUID, refreshToken string, accessExpAt, refreshExpAt time.Time) *Session {
	return &Session{
		id:           id,
		userID:       userID,
		refreshToken: refreshToken,
		accessExpAt:  accessExpAt,
		refreshExpAt: refreshExpAt,
	}
}

// Getter methods
func (s *Session) GetID() common.UUID {
	return s.id
}

func (s *Session) GetUserID() common.UUID {
	return s.userID
}

func (s *Session) GetRefreshToken() string {
	return s.refreshToken
}

func (s *Session) GetAccessExpAt() time.Time {
	return s.accessExpAt
}

func (s *Session) GetRefreshExpAt() time.Time {
	return s.refreshExpAt
}
