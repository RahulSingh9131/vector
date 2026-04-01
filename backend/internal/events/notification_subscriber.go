package events

import (
	"context"
	"encoding/json"
	"fmt"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// NotificationSubscriber is an Asynq task handler that receives domain events
// and creates notifications for relevant users.
type NotificationSubscriber struct {
	server           *server.Server
	notificationRepo *repository.NotificationRepository
	issueRepo        *repository.IssueRepository
	projectRepo      *repository.ProjectRepository
	logger           *zerolog.Logger
	// manager *realtime.NotificationManager // This will be added via Server struct
}

// NewNotificationSubscriber creates a new notification subscriber.
func NewNotificationSubscriber(s *server.Server, repos *repository.Repositories, logger *zerolog.Logger) *NotificationSubscriber {
	return &NotificationSubscriber{
		server:           s,
		notificationRepo: repos.Notification,
		issueRepo:        repos.Issue,
		projectRepo:      repos.Project,
		logger:           logger,
	}
}

// HandleEvent is the unified Asynq task handler for generating notifications.
func (s *NotificationSubscriber) HandleEvent(ctx context.Context, task *asynq.Task) error {
	var event DomainEvent
	if err := json.Unmarshal(task.Payload(), &event); err != nil {
		return fmt.Errorf("failed to unmarshal domain event for notification: %w", err)
	}

	notifications, err := s.mapToNotifications(ctx, event)
	if err != nil {
		s.logger.Error().Err(err).Str("event_type", event.Type).Msg("failed to map event to notifications")
		return err
	}

	for _, params := range notifications {
		// 1. Persist to DB
		notification, err := s.notificationRepo.Create(ctx, params)
		if err != nil {
			s.logger.Error().Err(err).Str("user_id", params.UserID.String()).Msg("failed to persist notification")
			continue
		}

		// 2. Broadcast to real-time manager
		if err := s.server.NotificationManager.Publish(ctx, notification); err != nil {
			s.logger.Error().Err(err).Str("notification_id", notification.ID.String()).Msg("failed to broadcast notification")
		}
	}

	return nil
}

// mapToNotifications determines which users should be notified based on the event.
func (s *NotificationSubscriber) mapToNotifications(ctx context.Context, event DomainEvent) ([]models.CreateNotificationParams, error) {
	var results []models.CreateNotificationParams

	switch event.Type {
	case "issue.assigned":
		var p IssueAssignedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}

		if p.NewAssigneeID != nil && *p.NewAssigneeID != event.ActorID {
			results = append(results, models.CreateNotificationParams{
				UserID:    *p.NewAssigneeID,
				ActorID:   &event.ActorID,
				ProjectID: &event.ProjectID,
				IssueID:   &p.IssueID,
				Type:      "issue_assigned",
				Title:     "New Issue Assigned",
				Message:   fmt.Sprintf("You have been assigned to issue %s", p.IssueID.String()),
				Payload:   event.Payload,
			})
		}

	case "comment.created":
		var p CommentCreatedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}

		// Notify issue assignee (if not the commenter)
		issue, err := s.issueRepo.GetByID(ctx, p.IssueID)
		if err != nil {
			return nil, err
		}
		if issue != nil && issue.AssigneeID != nil && *issue.AssigneeID != event.ActorID {
			results = append(results, models.CreateNotificationParams{
				UserID:    *issue.AssigneeID,
				ActorID:   &event.ActorID,
				ProjectID: &event.ProjectID,
				IssueID:   &p.IssueID,
				Type:      "comment_added",
				Title:     "New Comment",
				Message:   fmt.Sprintf("New comment on issue %s: %s", issue.ID.String(), p.BodyPreview),
				Payload:   event.Payload,
			})
		}

	case "member.added":
		var p MemberAddedPayload
		if err := json.Unmarshal(event.Payload, &p); err != nil {
			return nil, err
		}

		if p.UserID != event.ActorID {
			results = append(results, models.CreateNotificationParams{
				UserID:    p.UserID,
				ActorID:   &event.ActorID,
				ProjectID: &event.ProjectID,
				Type:      "project_joined",
				Title:     "Welcome to Project",
				Message:   fmt.Sprintf("You have been added to the project as a %s", p.Role),
				Payload:   event.Payload,
			})
		}
	}

	return results, nil
}
