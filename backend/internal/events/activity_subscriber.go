package events

import (
	"context"
	"encoding/json"
	"fmt"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// ActivitySubscriber is an Asynq task handler that receives domain events
// and writes corresponding records to the activities table.
type ActivitySubscriber struct {
	activityRepo *repository.ActivityRepository
	logger       *zerolog.Logger
}

// NewActivitySubscriber creates a new activity subscriber.
func NewActivitySubscriber(activityRepo *repository.ActivityRepository, logger *zerolog.Logger) *ActivitySubscriber {
	return &ActivitySubscriber{
		activityRepo: activityRepo,
		logger:       logger,
	}
}

// HandleEvent is the unified Asynq task handler for all domain events.
// It deserializes the DomainEvent envelope, maps it to a CreateActivityParams,
// and inserts the activity record.
func (s *ActivitySubscriber) HandleEvent(ctx context.Context, task *asynq.Task) error {
	var event DomainEvent
	if err := json.Unmarshal(task.Payload(), &event); err != nil {
		return fmt.Errorf("failed to unmarshal domain event: %w", err)
	}

	params, err := s.mapToActivity(event)
	if err != nil {
		s.logger.Error().Err(err).Str("event_type", event.Type).Msg("failed to map event to activity")
		return err
	}

	if _, err := s.activityRepo.Create(ctx, *params); err != nil {
		return fmt.Errorf("failed to create activity for event %s: %w", event.Type, err)
	}

	s.logger.Debug().
		Str("event_type", event.Type).
		Str("project_id", event.ProjectID.String()).
		Msg("activity recorded from event")

	return nil
}

// mapToActivity converts a DomainEvent into CreateActivityParams.
func (s *ActivitySubscriber) mapToActivity(event DomainEvent) (*models.CreateActivityParams, error) {
	base := models.CreateActivityParams{
		ProjectID: event.ProjectID,
		ActorID:   event.ActorID,
		Action:    event.Type,
	}

	switch event.Type {
	case "issue.created":
		var p IssueCreatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "issue"
		base.EntityID = p.IssueID
		base.NewValue = map[string]interface{}{
			"title":    p.Title,
			"status":   p.Status,
			"priority": p.Priority,
			"type":     p.Type,
		}

	case "issue.updated":
		var p IssueUpdatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "issue"
		base.EntityID = p.IssueID
		base.OldValue = p.OldValue
		base.NewValue = p.NewValue

	case "issue.assigned", "issue.unassigned":
		var p IssueAssignedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "issue"
		base.EntityID = p.IssueID
		base.OldValue = map[string]interface{}{"assignee_id": uuidPtrToStr(p.OldAssigneeID)}
		base.NewValue = map[string]interface{}{"assignee_id": uuidPtrToStr(p.NewAssigneeID)}

	case "issue.status_changed":
		var p IssueStatusChangedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "issue"
		base.EntityID = p.IssueID
		base.OldValue = map[string]interface{}{"status": p.OldStatus}
		base.NewValue = map[string]interface{}{"status": p.NewStatus}

	case "issue.deleted":
		var p IssueDeletedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		// issue_id is nil for deleted issues
		base.EntityType = "issue"
		base.EntityID = p.IssueID
		base.Metadata = map[string]interface{}{"title": p.Title, "issue_key": p.IssueKey}

	case "comment.created":
		var p CommentCreatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "comment"
		base.EntityID = p.CommentID
		base.Metadata = map[string]interface{}{"body_preview": p.BodyPreview}

	case "comment.updated":
		var p CommentUpdatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "comment"
		base.EntityID = p.CommentID

	case "comment.deleted":
		var p CommentDeletedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "comment"
		base.EntityID = p.CommentID

	case "label.added":
		var p LabelAddedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "label"
		base.EntityID = p.LabelID
		base.NewValue = map[string]interface{}{"label_name": p.LabelName, "color": p.Color}

	case "label.removed":
		var p LabelRemovedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.IssueID = &p.IssueID
		base.EntityType = "label"
		base.EntityID = p.LabelID
		base.OldValue = map[string]interface{}{"label_name": p.LabelName, "color": p.Color}

	case "project.created":
		var p ProjectCreatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "project"
		base.EntityID = p.ProjectID
		base.NewValue = map[string]interface{}{
			"name":       p.Name,
			"identifier": p.Identifier,
		}

	case "project.updated":
		var p ProjectUpdatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "project"
		base.EntityID = p.ProjectID
		base.OldValue = p.OldValue
		base.NewValue = p.NewValue

	case "project.deleted":
		var p ProjectDeletedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "project"
		base.EntityID = p.ProjectID
		base.Metadata = map[string]interface{}{"name": p.Name, "identifier": p.Identifier}

	case "member.added":
		var p MemberAddedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "member"
		base.EntityID = p.UserID
		base.NewValue = map[string]interface{}{"role": p.Role}

	case "member.removed":
		var p MemberRemovedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "member"
		base.EntityID = p.UserID

	case "member.role_changed":
		var p MemberRoleChangedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}
		base.EntityType = "member"
		base.EntityID = p.UserID
		base.OldValue = map[string]interface{}{"role": p.OldRole}
		base.NewValue = map[string]interface{}{"role": p.NewRole}

	default:
		return nil, fmt.Errorf("unknown event type: %s", event.Type)
	}

	return &base, nil
}

// uuidPtrToStr converts a *uuid.UUID to a string or nil for JSON serialization.
func uuidPtrToStr(id *uuid.UUID) interface{} {
	if id == nil {
		return nil
	}
	return id.String()
}
