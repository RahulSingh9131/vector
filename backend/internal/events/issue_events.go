package events

import "github.com/google/uuid"

// --- Issue Events ---

// IssueCreatedPayload is published when a new issue is created.
type IssueCreatedPayload struct {
	IssueID  uuid.UUID `json:"issue_id"`
	IssueKey string    `json:"issue_key"`
	Title    string    `json:"title"`
	Status   string    `json:"status"`
	Priority string    `json:"priority"`
	Type     string    `json:"type"`
}

// IssueUpdatedPayload is published when issue fields are modified.
type IssueUpdatedPayload struct {
	IssueID  uuid.UUID              `json:"issue_id"`
	OldValue map[string]interface{} `json:"old_value"`
	NewValue map[string]interface{} `json:"new_value"`
}

// IssueAssignedPayload is published when an issue is assigned or unassigned.
type IssueAssignedPayload struct {
	IssueID       uuid.UUID  `json:"issue_id"`
	OldAssigneeID *uuid.UUID `json:"old_assignee_id"`
	NewAssigneeID *uuid.UUID `json:"new_assignee_id"`
}

// IssueStatusChangedPayload is published when an issue's status changes.
type IssueStatusChangedPayload struct {
	IssueID   uuid.UUID `json:"issue_id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
}

// IssueDeletedPayload is published when an issue is deleted.
type IssueDeletedPayload struct {
	IssueID  uuid.UUID `json:"issue_id"`
	IssueKey string    `json:"issue_key"`
	Title    string    `json:"title"`
}
