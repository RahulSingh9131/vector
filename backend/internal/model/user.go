package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ClerkUserID  string     `json:"clerkUserId" db:"clerk_user_id"`
	Email        string     `json:"email" db:"email"`
	FirstName    *string    `json:"firstName,omitempty" db:"first_name"`
	LastName     *string    `json:"lastName,omitempty" db:"last_name"`
	AvatarURL    *string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	IsActive     bool       `json:"isActive" db:"is_active"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" db:"last_login_at"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// CreateUserParams represents the parameters for creating a new user
type CreateUserParams struct {
	ClerkUserID string
	Email       string
	FirstName   *string
	LastName    *string
	AvatarURL   *string
}

// UpdateUserParams represents the parameters for updating a user
type UpdateUserParams struct {
	Email     *string
	FirstName *string
	LastName  *string
	AvatarURL *string
	IsActive  *bool
}
