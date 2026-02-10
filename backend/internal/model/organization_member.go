package models

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationMember represents a user's membership in an organization
type OrganizationMember struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organizationId" db:"organization_id"`
	UserID         uuid.UUID `json:"userId" db:"user_id"`
	Role           string    `json:"role" db:"role"`
	JoinedAt       time.Time `json:"joinedAt" db:"joined_at"`
}

// OrganizationMemberWithDetails includes user and organization details
type OrganizationMemberWithDetails struct {
	OrganizationMember
	User         *User         `json:"user,omitempty"`
	Organization *Organization `json:"organization,omitempty"`
}

// CreateOrganizationMemberParams represents the parameters for adding a member to an organization
type CreateOrganizationMemberParams struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
}
