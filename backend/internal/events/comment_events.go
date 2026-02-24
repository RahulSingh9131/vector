package events

import "github.com/google/uuid"

// --- Comment Events ---

// CommentCreatedPayload is published when a new comment is added.
type CommentCreatedPayload struct {
	CommentID   uuid.UUID `json:"comment_id"`
	IssueID     uuid.UUID `json:"issue_id"`
	BodyPreview string    `json:"body_preview"`
}

// CommentUpdatedPayload is published when a comment is edited.
type CommentUpdatedPayload struct {
	CommentID uuid.UUID `json:"comment_id"`
	IssueID   uuid.UUID `json:"issue_id"`
}

// CommentDeletedPayload is published when a comment is soft-deleted.
type CommentDeletedPayload struct {
	CommentID uuid.UUID `json:"comment_id"`
	IssueID   uuid.UUID `json:"issue_id"`
}
