package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notification represents a notification sent to a user.
type Notification struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	ActorID   *uuid.UUID      `json:"actor_id" db:"actor_id"`
	ProjectID *uuid.UUID      `json:"project_id" db:"project_id"`
	IssueID   *uuid.UUID      `json:"issue_id" db:"issue_id"`
	Type      string          `json:"type" db:"type"`
	Title     string          `json:"title" db:"title"`
	Message   string          `json:"message" db:"message"`
	Payload   json.RawMessage `json:"payload" db:"payload"`
	IsRead    bool            `json:"is_read" db:"is_read"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// CreateNotificationParams defines the input for creating a new notification.
type CreateNotificationParams struct {
	UserID    uuid.UUID       `json:"user_id"`
	ActorID   *uuid.UUID      `json:"actor_id"`
	ProjectID *uuid.UUID      `json:"project_id"`
	IssueID   *uuid.UUID      `json:"issue_id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	Payload   json.RawMessage `json:"payload"`
}

// NotificationFilters defines filtering options for listing notifications.
type NotificationFilters struct {
	UserID *uuid.UUID `json:"user_id"`
	IsRead *bool      `json:"is_read"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
