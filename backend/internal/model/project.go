package models

import (
	"time"

	"github.com/google/uuid"
)

// Project represents a project within an organization
type Project struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organizationId" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Slug           string    `json:"slug" db:"slug"`
	Description    *string   `json:"description,omitempty" db:"description"`
	Status         string    `json:"status" db:"status"`
	Identifier     string    `json:"identifier" db:"identifier"`
	IssueCounter   int       `json:"issueCounter" db:"issue_counter"`
	CreatedBy      uuid.UUID `json:"createdBy" db:"created_by"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateProjectParams represents the parameters for creating a new project
type CreateProjectParams struct {
	OrganizationID uuid.UUID
	Name           string
	Slug           string
	Description    *string
	Identifier     string
	CreatedBy      uuid.UUID
}

// UpdateProjectParams represents the parameters for updating a project
type UpdateProjectParams struct {
	Name        *string
	Slug        *string
	Description *string
	Status      *string
}
