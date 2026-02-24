package events

import "github.com/google/uuid"

// --- Label Events ---

// LabelAddedPayload is published when a label is attached to an issue.
type LabelAddedPayload struct {
	IssueID   uuid.UUID `json:"issue_id"`
	LabelID   uuid.UUID `json:"label_id"`
	LabelName string    `json:"label_name"`
	Color     string    `json:"color"`
}

// LabelRemovedPayload is published when a label is detached from an issue.
type LabelRemovedPayload struct {
	IssueID   uuid.UUID `json:"issue_id"`
	LabelID   uuid.UUID `json:"label_id"`
	LabelName string    `json:"label_name"`
	Color     string    `json:"color"`
}
