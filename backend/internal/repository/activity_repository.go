// Package repository provides database access layer implementations
// for managing application entities such as activities, issues, and projects.
package repository

import (
	"context"
	"encoding/json"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

// ActivityRepository handles database operations for activities
type ActivityRepository struct {
	server *server.Server
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(s *server.Server) *ActivityRepository {
	return &ActivityRepository{
		server: s,
	}
}

// Create records a new activity entry
func (r *ActivityRepository) Create(ctx context.Context, params models.CreateActivityParams) (*models.Activity, error) {
	// Marshal JSONB fields
	var oldValueJSON, newValueJSON, metadataJSON []byte
	var err error

	if params.OldValue != nil {
		oldValueJSON, err = json.Marshal(params.OldValue)
		if err != nil {
			return nil, err
		}
	}
	if params.NewValue != nil {
		newValueJSON, err = json.Marshal(params.NewValue)
		if err != nil {
			return nil, err
		}
	}
	if params.Metadata != nil {
		metadataJSON, err = json.Marshal(params.Metadata)
		if err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO activities (project_id, issue_id, actor_id, action, entity_type, entity_id, old_value, new_value, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, project_id, issue_id, actor_id, action, entity_type, entity_id,
		          old_value, new_value, metadata, created_at
	`

	var activity models.Activity
	err = r.server.DB.Pool.QueryRow(
		ctx, query,
		params.ProjectID, params.IssueID, params.ActorID,
		params.Action, params.EntityType, params.EntityID,
		nullableJSON(oldValueJSON), nullableJSON(newValueJSON), nullableJSON(metadataJSON),
	).Scan(
		&activity.ID, &activity.ProjectID, &activity.IssueID, &activity.ActorID,
		&activity.Action, &activity.EntityType, &activity.EntityID,
		&activity.OldValue, &activity.NewValue, &activity.Metadata, &activity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

// ListByIssue retrieves activities for an issue (newest first, paginated)
func (r *ActivityRepository) ListByIssue(ctx context.Context, issueID uuid.UUID, page, limit int) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Count total
	countQuery := `SELECT COUNT(*) FROM activities WHERE issue_id = $1`
	var total int
	if err := r.server.DB.Pool.QueryRow(ctx, countQuery, issueID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := `
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		WHERE a.issue_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, issueID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []models.ActivityWithActor
	for rows.Next() {
		var a models.ActivityWithActor
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.IssueID, &a.ActorID,
			&a.Action, &a.EntityType, &a.EntityID,
			&a.OldValue, &a.NewValue, &a.Metadata, &a.CreatedAt,
			&a.ActorFirstName, &a.ActorLastName,
			&a.ActorAvatarURL, &a.ActorEmail,
		); err != nil {
			return nil, 0, err
		}
		activities = append(activities, a)
	}

	if activities == nil {
		activities = []models.ActivityWithActor{}
	}

	return activities, total, nil
}

// ListByProject retrieves activities for a project (newest first, paginated)
func (r *ActivityRepository) ListByProject(ctx context.Context, projectID uuid.UUID, page, limit int) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Count total
	countQuery := `SELECT COUNT(*) FROM activities WHERE project_id = $1`
	var total int
	if err := r.server.DB.Pool.QueryRow(ctx, countQuery, projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := `
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		WHERE a.project_id = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var activities []models.ActivityWithActor
	for rows.Next() {
		var a models.ActivityWithActor
		if err := rows.Scan(
			&a.ID, &a.ProjectID, &a.IssueID, &a.ActorID,
			&a.Action, &a.EntityType, &a.EntityID,
			&a.OldValue, &a.NewValue, &a.Metadata, &a.CreatedAt,
			&a.ActorFirstName, &a.ActorLastName,
			&a.ActorAvatarURL, &a.ActorEmail,
		); err != nil {
			return nil, 0, err
		}
		activities = append(activities, a)
	}

	if activities == nil {
		activities = []models.ActivityWithActor{}
	}

	return activities, total, nil
}

// nullableJSON returns nil if the byte slice is nil/empty, otherwise returns the slice
func nullableJSON(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	return b
}
