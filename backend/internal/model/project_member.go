package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectMember represents a user's membership in a project
type ProjectMember struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProjectID uuid.UUID `json:"projectId" db:"project_id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Role      string    `json:"role" db:"role"`
	JoinedAt  time.Time `json:"joinedAt" db:"joined_at"`
}

// ProjectMemberWithDetails includes user and project details
type ProjectMemberWithDetails struct {
	ProjectMember
	User    *User    `json:"user,omitempty"`
	Project *Project `json:"project,omitempty"`
}

// CreateProjectMemberParams represents the parameters for adding a member to a project
type CreateProjectMemberParams struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Role      string
}
