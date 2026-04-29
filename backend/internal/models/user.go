package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UID       uuid.UUID `json:"uid"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCredentials carries auth-only fields that must never be exposed via JSON.
// It is returned by repository lookups used by the auth handler.
type UserCredentials struct {
	UID          uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	IsAdmin      bool
}
