package repository

import (
	"context"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

// NotificationRepository handles database operations for notifications.
type NotificationRepository struct {
	server *server.Server
	db     database.DBTX
}

// NewNotificationRepository creates a new notification repository.
func NewNotificationRepository(s *server.Server, db database.DBTX) *NotificationRepository {
	return &NotificationRepository{
		server: s,
		db:     db,
	}
}

// Create inserts a new notification into the database.
func (r *NotificationRepository) Create(ctx context.Context, params models.CreateNotificationParams) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (
			user_id, actor_id, project_id, issue_id, type, title, message, payload
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, user_id, actor_id, project_id, issue_id, type, title, message, payload, is_read, created_at
	`

	var n models.Notification
	err := r.db.QueryRow(ctx, query,
		params.UserID,
		params.ActorID,
		params.ProjectID,
		params.IssueID,
		params.Type,
		params.Title,
		params.Message,
		params.Payload,
	).Scan(
		&n.ID,
		&n.UserID,
		&n.ActorID,
		&n.ProjectID,
		&n.IssueID,
		&n.Type,
		&n.Title,
		&n.Message,
		&n.Payload,
		&n.IsRead,
		&n.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &n, nil
}

// GetByID retrieves a single notification by its ID.
func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, project_id, issue_id, type, title, message, payload, is_read, created_at
		FROM notifications
		WHERE id = $1
	`

	var n models.Notification
	err := r.db.QueryRow(ctx, query, id).Scan(
		&n.ID,
		&n.UserID,
		&n.ActorID,
		&n.ProjectID,
		&n.IssueID,
		&n.Type,
		&n.Title,
		&n.Message,
		&n.Payload,
		&n.IsRead,
		&n.CreatedAt,
	)

	if err != nil {
		return nil, nil // Return nil for not found (matching existing repo patterns)
	}

	return &n, nil
}

// ListByUser retrieves notifications for a specific user with filtering and pagination.
func (r *NotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, filters models.NotificationFilters) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, actor_id, project_id, issue_id, type, title, message, payload, is_read, created_at
		FROM notifications
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIdx := 2

	if filters.IsRead != nil {
		query += fmt.Sprintf(" AND is_read = $%d", argIdx)
		args = append(args, *filters.IsRead)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filters.Limit)
		argIdx++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filters.Offset)
		argIdx++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.ActorID,
			&n.ProjectID,
			&n.IssueID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Payload,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkAsRead marks a specific notification as read.
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// MarkAllAsRead marks all notifications for a user as read.
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// Delete removes a notification from the database.
func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
