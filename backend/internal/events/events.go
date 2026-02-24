// Package events provides domain event types, an event bus abstraction,
// and subscribers for asynchronous activity recording.
package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the envelope that wraps all domain events.
// It contains routing metadata and a JSON-encoded payload.
type DomainEvent struct {
	// Type identifies the event, e.g. "issue.created", "comment.deleted".
	Type string `json:"type"`

	// ProjectID scopes the event to a project (for activity recording).
	ProjectID uuid.UUID `json:"project_id"`

	// ActorID is the user who triggered this event.
	ActorID uuid.UUID `json:"actor_id"`

	// Payload is the JSON-encoded, event-specific data.
	Payload json.RawMessage `json:"payload"`

	// Timestamp records when the event occurred.
	Timestamp time.Time `json:"timestamp"`
}
