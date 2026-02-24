package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// EventPublisher provides typed helper methods for publishing domain events.
// Services use this instead of constructing DomainEvent structs directly,
// ensuring consistent event structure and reducing boilerplate.
type EventPublisher struct {
	bus    EventBus
	logger *zerolog.Logger
}

// NewEventPublisher creates a new event publisher.
func NewEventPublisher(bus EventBus, logger *zerolog.Logger) *EventPublisher {
	return &EventPublisher{
		bus:    bus,
		logger: logger,
	}
}

// publish is the internal helper that builds a DomainEvent and publishes it.
func (p *EventPublisher) publish(ctx context.Context, eventType string, projectID, actorID uuid.UUID, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error().Err(err).Str("event_type", eventType).Msg("failed to marshal event payload")
		return
	}

	event := DomainEvent{
		Type:      eventType,
		ProjectID: projectID,
		ActorID:   actorID,
		Payload:   data,
		Timestamp: time.Now(),
	}

	if err := p.bus.Publish(ctx, event); err != nil {
		p.logger.Error().Err(err).Str("event_type", eventType).Msg("failed to publish event")
	}
}

// --- Issue Events ---

// IssueCreated publishes an issue.created event.
func (p *EventPublisher) IssueCreated(ctx context.Context, projectID, actorID uuid.UUID, payload IssueCreatedPayload) {
	p.publish(ctx, "issue.created", projectID, actorID, payload)
}

// IssueUpdated publishes an issue.updated event.
func (p *EventPublisher) IssueUpdated(ctx context.Context, projectID, actorID uuid.UUID, payload IssueUpdatedPayload) {
	p.publish(ctx, "issue.updated", projectID, actorID, payload)
}

// IssueAssigned publishes an issue.assigned or issue.unassigned event.
func (p *EventPublisher) IssueAssigned(ctx context.Context, projectID, actorID uuid.UUID, payload IssueAssignedPayload) {
	action := "issue.assigned"
	if payload.NewAssigneeID == nil {
		action = "issue.unassigned"
	}
	p.publish(ctx, action, projectID, actorID, payload)
}

// IssueStatusChanged publishes an issue.status_changed event.
func (p *EventPublisher) IssueStatusChanged(ctx context.Context, projectID, actorID uuid.UUID, payload IssueStatusChangedPayload) {
	p.publish(ctx, "issue.status_changed", projectID, actorID, payload)
}

// IssueDeleted publishes an issue.deleted event.
func (p *EventPublisher) IssueDeleted(ctx context.Context, projectID, actorID uuid.UUID, payload IssueDeletedPayload) {
	p.publish(ctx, "issue.deleted", projectID, actorID, payload)
}

// --- Comment Events ---

// CommentCreated publishes a comment.created event.
func (p *EventPublisher) CommentCreated(ctx context.Context, projectID, actorID uuid.UUID, payload CommentCreatedPayload) {
	p.publish(ctx, "comment.created", projectID, actorID, payload)
}

// CommentUpdated publishes a comment.updated event.
func (p *EventPublisher) CommentUpdated(ctx context.Context, projectID, actorID uuid.UUID, payload CommentUpdatedPayload) {
	p.publish(ctx, "comment.updated", projectID, actorID, payload)
}

// CommentDeleted publishes a comment.deleted event.
func (p *EventPublisher) CommentDeleted(ctx context.Context, projectID, actorID uuid.UUID, payload CommentDeletedPayload) {
	p.publish(ctx, "comment.deleted", projectID, actorID, payload)
}

// --- Label Events ---

// LabelAdded publishes a label.added event.
func (p *EventPublisher) LabelAdded(ctx context.Context, projectID, actorID uuid.UUID, payload LabelAddedPayload) {
	p.publish(ctx, "label.added", projectID, actorID, payload)
}

// LabelRemoved publishes a label.removed event.
func (p *EventPublisher) LabelRemoved(ctx context.Context, projectID, actorID uuid.UUID, payload LabelRemovedPayload) {
	p.publish(ctx, "label.removed", projectID, actorID, payload)
}

// --- Project Events ---

// ProjectCreated publishes a project.created event.
func (p *EventPublisher) ProjectCreated(ctx context.Context, projectID, actorID uuid.UUID, payload ProjectCreatedPayload) {
	p.publish(ctx, "project.created", projectID, actorID, payload)
}

// ProjectUpdated publishes a project.updated event.
func (p *EventPublisher) ProjectUpdated(ctx context.Context, projectID, actorID uuid.UUID, payload ProjectUpdatedPayload) {
	p.publish(ctx, "project.updated", projectID, actorID, payload)
}

// ProjectDeleted publishes a project.deleted event.
func (p *EventPublisher) ProjectDeleted(ctx context.Context, projectID, actorID uuid.UUID, payload ProjectDeletedPayload) {
	p.publish(ctx, "project.deleted", projectID, actorID, payload)
}

// --- Member Events ---

// MemberAdded publishes a member.added event.
func (p *EventPublisher) MemberAdded(ctx context.Context, projectID, actorID uuid.UUID, payload MemberAddedPayload) {
	p.publish(ctx, "member.added", projectID, actorID, payload)
}

// MemberRemoved publishes a member.removed event.
func (p *EventPublisher) MemberRemoved(ctx context.Context, projectID, actorID uuid.UUID, payload MemberRemovedPayload) {
	p.publish(ctx, "member.removed", projectID, actorID, payload)
}

// MemberRoleChanged publishes a member.role_changed event.
func (p *EventPublisher) MemberRoleChanged(ctx context.Context, projectID, actorID uuid.UUID, payload MemberRoleChangedPayload) {
	p.publish(ctx, "member.role_changed", projectID, actorID, payload)
}
