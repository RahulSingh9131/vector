package models

import (
	"time"

	"github.com/google/uuid"
)

// Issue represents a task, bug, story, or epic within a project
type Issue struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ProjectID     uuid.UUID  `json:"projectId" db:"project_id"`
	IssueKey      string     `json:"issueKey" db:"issue_key"`
	Title         string     `json:"title" db:"title"`
	Description   *string    `json:"description,omitempty" db:"description"`
	Status        string     `json:"status" db:"status"`
	Priority      string     `json:"priority" db:"priority"`
	Type          string     `json:"type" db:"type"`
	AssigneeID    *uuid.UUID `json:"assigneeId,omitempty" db:"assignee_id"`
	ReporterID    uuid.UUID  `json:"reporterId" db:"reporter_id"`
	SortOrder     int        `json:"sortOrder" db:"sort_order"`
	ParentIssueID *uuid.UUID `json:"parentIssueId,omitempty" db:"parent_issue_id"`
	DueDate       *time.Time `json:"dueDate,omitempty" db:"due_date"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

// IssueWithDetails includes assignee and reporter user details
type IssueWithDetails struct {
	Issue
	Assignee *User `json:"assignee,omitempty"`
	Reporter *User `json:"reporter,omitempty"`
}

// CreateIssueParams represents the parameters for creating a new issue
type CreateIssueParams struct {
	ProjectID     uuid.UUID
	Title         string
	Description   *string
	Priority      string
	Type          string
	AssigneeID    *uuid.UUID
	ReporterID    uuid.UUID
	ParentIssueID *uuid.UUID
	DueDate       *time.Time
}

// UpdateIssueParams represents the parameters for updating an issue
type UpdateIssueParams struct {
	Title         *string
	Description   *string
	Status        *string
	Priority      *string
	Type          *string
	AssigneeID    *uuid.UUID
	SortOrder     *int
	ParentIssueID *uuid.UUID
	DueDate       *time.Time
}

// IssueFilters represents filtering options for listing issues
type IssueFilters struct {
	Status     *string
	Priority   *string
	Type       *string
	AssigneeID *uuid.UUID
	Page       int
	Limit      int
}
