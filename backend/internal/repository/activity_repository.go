// Package repository provides database access layer implementations
// for managing application entities such as activities, issues, and projects.
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

// ActivityRepository handles database operations for activities
type ActivityRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewActivityRepository creates a new activity repository
func NewActivityRepository(s *server.Server) *ActivityRepository {
	return &ActivityRepository{
		db:     s.DB.Pool,
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
	err = r.db.QueryRow(
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

// ListByIssue retrieves activities for an issue (newest first, paginated, with optional filters)
func (r *ActivityRepository) ListByIssue(ctx context.Context, issueID uuid.UUID, page, limit int, filters models.ActivityFilters) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build dynamic WHERE clause
	where := "WHERE a.issue_id = $1"
	args := []interface{}{issueID}
	argIdx := 2

	if filters.Action != nil {
		where += fmt.Sprintf(" AND a.action = $%d", argIdx)
		args = append(args, *filters.Action)
		argIdx++
	}
	if filters.EntityType != nil {
		where += fmt.Sprintf(" AND a.entity_type = $%d", argIdx)
		args = append(args, *filters.EntityType)
		argIdx++
	}
	if filters.ActorID != nil {
		where += fmt.Sprintf(" AND a.actor_id = $%d", argIdx)
		args = append(args, *filters.ActorID)
		argIdx++
	}
	if filters.From != nil {
		where += fmt.Sprintf(" AND a.created_at >= $%d", argIdx)
		args = append(args, *filters.From)
		argIdx++
	}
	if filters.To != nil {
		where += fmt.Sprintf(" AND a.created_at <= $%d", argIdx)
		args = append(args, *filters.To)
		argIdx++
	}

	// Count total with filters
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM activities a %s", where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := fmt.Sprintf(`
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
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

// ListByProject retrieves activities for a project (newest first, paginated, with optional filters)
func (r *ActivityRepository) ListByProject(ctx context.Context, projectID uuid.UUID, page, limit int, filters models.ActivityFilters) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build dynamic WHERE clause
	where := "WHERE a.project_id = $1"
	args := []interface{}{projectID}
	argIdx := 2

	if filters.Action != nil {
		where += fmt.Sprintf(" AND a.action = $%d", argIdx)
		args = append(args, *filters.Action)
		argIdx++
	}
	if filters.EntityType != nil {
		where += fmt.Sprintf(" AND a.entity_type = $%d", argIdx)
		args = append(args, *filters.EntityType)
		argIdx++
	}
	if filters.ActorID != nil {
		where += fmt.Sprintf(" AND a.actor_id = $%d", argIdx)
		args = append(args, *filters.ActorID)
		argIdx++
	}
	if filters.From != nil {
		where += fmt.Sprintf(" AND a.created_at >= $%d", argIdx)
		args = append(args, *filters.From)
		argIdx++
	}
	if filters.To != nil {
		where += fmt.Sprintf(" AND a.created_at <= $%d", argIdx)
		args = append(args, *filters.To)
		argIdx++
	}

	// Count total with filters
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM activities a %s", where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := fmt.Sprintf(`
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
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

// ListByActor retrieves activities performed by a specific user across all projects (newest first, paginated, with optional filters)
func (r *ActivityRepository) ListByActor(ctx context.Context, actorID uuid.UUID, page, limit int, filters models.ActivityFilters) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build dynamic WHERE clause — actor_id is the primary constraint
	where := "WHERE a.actor_id = $1"
	args := []interface{}{actorID}
	argIdx := 2

	if filters.Action != nil {
		where += fmt.Sprintf(" AND a.action = $%d", argIdx)
		args = append(args, *filters.Action)
		argIdx++
	}
	if filters.EntityType != nil {
		where += fmt.Sprintf(" AND a.entity_type = $%d", argIdx)
		args = append(args, *filters.EntityType)
		argIdx++
	}
	if filters.From != nil {
		where += fmt.Sprintf(" AND a.created_at >= $%d", argIdx)
		args = append(args, *filters.From)
		argIdx++
	}
	if filters.To != nil {
		where += fmt.Sprintf(" AND a.created_at <= $%d", argIdx)
		args = append(args, *filters.To)
		argIdx++
	}

	// Count total with filters
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM activities a %s", where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := fmt.Sprintf(`
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
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

// ListByOrganization retrieves activities across all projects within an organization
// that the given user is a member of (newest first, paginated, with optional filters).
func (r *ActivityRepository) ListByOrganization(ctx context.Context, orgID, userID uuid.UUID, page, limit int, filters models.ActivityFilters) ([]models.ActivityWithActor, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Base WHERE: project belongs to org AND user is a member of the project
	where := "WHERE p.organization_id = $1 AND pm.user_id = $2"
	args := []interface{}{orgID, userID}
	argIdx := 3

	if filters.Action != nil {
		where += fmt.Sprintf(" AND a.action = $%d", argIdx)
		args = append(args, *filters.Action)
		argIdx++
	}
	if filters.EntityType != nil {
		where += fmt.Sprintf(" AND a.entity_type = $%d", argIdx)
		args = append(args, *filters.EntityType)
		argIdx++
	}
	if filters.ActorID != nil {
		where += fmt.Sprintf(" AND a.actor_id = $%d", argIdx)
		args = append(args, *filters.ActorID)
		argIdx++
	}
	if filters.From != nil {
		where += fmt.Sprintf(" AND a.created_at >= $%d", argIdx)
		args = append(args, *filters.From)
		argIdx++
	}
	if filters.To != nil {
		where += fmt.Sprintf(" AND a.created_at <= $%d", argIdx)
		args = append(args, *filters.To)
		argIdx++
	}

	// Count total with filters
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM activities a
		JOIN projects p ON p.id = a.project_id
		JOIN project_members pm ON pm.project_id = p.id
		%s
	`, where)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch with actor details
	query := fmt.Sprintf(`
		SELECT a.id, a.project_id, a.issue_id, a.actor_id, a.action, a.entity_type, a.entity_id,
		       a.old_value, a.new_value, a.metadata, a.created_at,
		       u.first_name AS actor_first_name, u.last_name AS actor_last_name,
		       u.avatar_url AS actor_avatar_url, u.email AS actor_email
		FROM activities a
		JOIN projects p ON p.id = a.project_id
		JOIN project_members pm ON pm.project_id = p.id
		LEFT JOIN users u ON u.id = a.actor_id
		%s
		ORDER BY a.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
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
