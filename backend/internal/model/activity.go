// Package models defines the domain data structures and types used throughout the application.
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Activity represents an activity log entry
type Activity struct {
	ID         uuid.UUID        `json:"id" db:"id"`
	ProjectID  uuid.UUID        `json:"project_id" db:"project_id"`
	IssueID    *uuid.UUID       `json:"issue_id" db:"issue_id"`
	ActorID    uuid.UUID        `json:"actor_id" db:"actor_id"`
	Action     string           `json:"action" db:"action"`
	EntityType string           `json:"entity_type" db:"entity_type"`
	EntityID   uuid.UUID        `json:"entity_id" db:"entity_id"`
	OldValue   *json.RawMessage `json:"old_value" db:"old_value"`
	NewValue   *json.RawMessage `json:"new_value" db:"new_value"`
	Metadata   *json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
}

// ActivityWithActor extends Activity with actor details
type ActivityWithActor struct {
	Activity
	ActorFirstName *string `json:"actor_first_name" db:"actor_first_name"`
	ActorLastName  *string `json:"actor_last_name" db:"actor_last_name"`
	ActorAvatarURL *string `json:"actor_avatar_url" db:"actor_avatar_url"`
	ActorEmail     string  `json:"actor_email" db:"actor_email"`
}

// CreateActivityParams contains the fields needed to record an activity
type CreateActivityParams struct {
	ProjectID  uuid.UUID
	IssueID    *uuid.UUID
	ActorID    uuid.UUID
	Action     string
	EntityType string
	EntityID   uuid.UUID
	OldValue   interface{}
	NewValue   interface{}
	Metadata   interface{}
}

// ActivityFilters holds optional filter criteria for listing activities
type ActivityFilters struct {
	Action     *string    // e.g. "issue.created"
	EntityType *string    // e.g. "issue", "comment", "label"
	ActorID    *uuid.UUID // filter by who performed the action
	From       *time.Time // activities created on or after this time
	To         *time.Time // activities created on or before this time
}
