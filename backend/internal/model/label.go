package models

import (
	"time"

	"github.com/google/uuid"
)

// Label represents a project-scoped label for categorizing issues
type Label struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateLabelParams contains the fields needed to create a label
type CreateLabelParams struct {
	ProjectID   uuid.UUID `json:"project_id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description *string   `json:"description,omitempty"`
}

// UpdateLabelParams contains the fields that can be updated on a label
type UpdateLabelParams struct {
	Name        *string `json:"name,omitempty"`
	Color       *string `json:"color,omitempty"`
	Description *string `json:"description,omitempty"`
}

// IssueLabel represents the association between an issue and a label
type IssueLabel struct {
	IssueID   uuid.UUID `json:"issue_id" db:"issue_id"`
	LabelID   uuid.UUID `json:"label_id" db:"label_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
