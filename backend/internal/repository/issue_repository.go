package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// IssueRepository handles database operations for issues
type IssueRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewIssueRepository creates a new issue repository
func NewIssueRepository(s *server.Server) *IssueRepository {
	return &IssueRepository{
		db:     s.DB.Pool,
		server: s,
	}
}

// Create creates a new issue
func (r *IssueRepository) Create(ctx context.Context, params models.CreateIssueParams, issueKey string) (*models.Issue, error) {
	query := `
		INSERT INTO issues (project_id, issue_key, title, description, priority, type,
		                    assignee_id, reporter_id, parent_issue_id, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, project_id, issue_key, title, description, status, priority, type,
		          assignee_id, reporter_id, sort_order, parent_issue_id, due_date,
		          created_at, updated_at
	`

	var issue models.Issue
	err := r.db.QueryRow(
		ctx,
		query,
		params.ProjectID,
		issueKey,
		params.Title,
		params.Description,
		params.Priority,
		params.Type,
		params.AssigneeID,
		params.ReporterID,
		params.ParentIssueID,
		params.DueDate,
	).Scan(
		&issue.ID,
		&issue.ProjectID,
		&issue.IssueKey,
		&issue.Title,
		&issue.Description,
		&issue.Status,
		&issue.Priority,
		&issue.Type,
		&issue.AssigneeID,
		&issue.ReporterID,
		&issue.SortOrder,
		&issue.ParentIssueID,
		&issue.DueDate,
		&issue.CreatedAt,
		&issue.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &issue, nil
}

// GetByID retrieves an issue by its UUID with user details
func (r *IssueRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.IssueWithDetails, error) {
	query := `
		SELECT i.id, i.project_id, i.issue_key, i.title, i.description, i.status,
		       i.priority, i.type, i.assignee_id, i.reporter_id, i.sort_order,
		       i.parent_issue_id, i.due_date, i.created_at, i.updated_at,
		       a.id, a.clerk_user_id, a.email, a.first_name, a.last_name, a.avatar_url,
		       a.is_active, a.last_login_at, a.created_at, a.updated_at,
		       r.id, r.clerk_user_id, r.email, r.first_name, r.last_name, r.avatar_url,
		       r.is_active, r.last_login_at, r.created_at, r.updated_at
		FROM issues i
		LEFT JOIN users a ON a.id = i.assignee_id
		INNER JOIN users r ON r.id = i.reporter_id
		WHERE i.id = $1
	`

	var issue models.IssueWithDetails
	var reporter models.User

	// Assignee fields are nullable
	var aID *uuid.UUID
	var aClerkID, aEmail *string
	var aFirstName, aLastName, aAvatarURL *string
	var aIsActive *bool
	var aLastLogin, aCreatedAt, aUpdatedAt *interface{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&issue.ID, &issue.ProjectID, &issue.IssueKey, &issue.Title, &issue.Description,
		&issue.Status, &issue.Priority, &issue.Type, &issue.AssigneeID, &issue.ReporterID,
		&issue.SortOrder, &issue.ParentIssueID, &issue.DueDate, &issue.CreatedAt, &issue.UpdatedAt,
		// Assignee (nullable)
		&aID, &aClerkID, &aEmail, &aFirstName, &aLastName, &aAvatarURL,
		&aIsActive, &aLastLogin, &aCreatedAt, &aUpdatedAt,
		// Reporter
		&reporter.ID, &reporter.ClerkUserID, &reporter.Email, &reporter.FirstName,
		&reporter.LastName, &reporter.AvatarURL, &reporter.IsActive, &reporter.LastLoginAt,
		&reporter.CreatedAt, &reporter.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	issue.Reporter = &reporter
	if aID != nil {
		issue.Assignee = &models.User{
			ID:          *aID,
			ClerkUserID: *aClerkID,
			Email:       *aEmail,
			FirstName:   aFirstName,
			LastName:    aLastName,
			AvatarURL:   aAvatarURL,
		}
	}

	return &issue, nil
}

// GetByIssueKey retrieves an issue by project ID and issue key
func (r *IssueRepository) GetByIssueKey(ctx context.Context, projectID uuid.UUID, issueKey string) (*models.Issue, error) {
	query := `
		SELECT id, project_id, issue_key, title, description, status, priority, type,
		       assignee_id, reporter_id, sort_order, parent_issue_id, due_date,
		       created_at, updated_at
		FROM issues
		WHERE project_id = $1 AND issue_key = $2
	`

	var issue models.Issue
	err := r.db.QueryRow(ctx, query, projectID, issueKey).Scan(
		&issue.ID, &issue.ProjectID, &issue.IssueKey, &issue.Title, &issue.Description,
		&issue.Status, &issue.Priority, &issue.Type, &issue.AssigneeID, &issue.ReporterID,
		&issue.SortOrder, &issue.ParentIssueID, &issue.DueDate, &issue.CreatedAt, &issue.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &issue, nil
}

// ListByProject retrieves issues for a project with optional filters and pagination
func (r *IssueRepository) ListByProject(ctx context.Context, projectID uuid.UUID, filters models.IssueFilters) (*models.PaginatedResponse[models.IssueWithDetails], error) {
	// Build WHERE clauses dynamically
	conditions := []string{"i.project_id = $1"}
	args := []interface{}{projectID}
	argIndex := 2

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("i.status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}
	if filters.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("i.priority = $%d", argIndex))
		args = append(args, *filters.Priority)
		argIndex++
	}
	if filters.Type != nil {
		conditions = append(conditions, fmt.Sprintf("i.type = $%d", argIndex))
		args = append(args, *filters.Type)
		argIndex++
	}
	if filters.AssigneeID != nil {
		conditions = append(conditions, fmt.Sprintf("i.assignee_id = $%d", argIndex))
		args = append(args, *filters.AssigneeID)
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM issues i WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Set defaults
	page := filters.Page
	if page < 1 {
		page = 1
	}
	limit := filters.Limit
	if limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit
	totalPages := (total + limit - 1) / limit

	// Fetch data
	dataQuery := fmt.Sprintf(`
		SELECT i.id, i.project_id, i.issue_key, i.title, i.description, i.status,
		       i.priority, i.type, i.assignee_id, i.reporter_id, i.sort_order,
		       i.parent_issue_id, i.due_date, i.created_at, i.updated_at,
		       r.id, r.clerk_user_id, r.email, r.first_name, r.last_name, r.avatar_url,
		       r.is_active, r.last_login_at, r.created_at, r.updated_at
		FROM issues i
		INNER JOIN users r ON r.id = i.reporter_id
		WHERE %s
		ORDER BY i.sort_order ASC, i.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.IssueWithDetails
	for rows.Next() {
		var issue models.IssueWithDetails
		var reporter models.User
		if err := rows.Scan(
			&issue.ID, &issue.ProjectID, &issue.IssueKey, &issue.Title, &issue.Description,
			&issue.Status, &issue.Priority, &issue.Type, &issue.AssigneeID, &issue.ReporterID,
			&issue.SortOrder, &issue.ParentIssueID, &issue.DueDate, &issue.CreatedAt, &issue.UpdatedAt,
			&reporter.ID, &reporter.ClerkUserID, &reporter.Email, &reporter.FirstName,
			&reporter.LastName, &reporter.AvatarURL, &reporter.IsActive, &reporter.LastLoginAt,
			&reporter.CreatedAt, &reporter.UpdatedAt,
		); err != nil {
			return nil, err
		}
		issue.Reporter = &reporter
		issues = append(issues, issue)
	}

	return &models.PaginatedResponse[models.IssueWithDetails]{
		Data:       issues,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// Update updates an issue
func (r *IssueRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateIssueParams) (*models.Issue, error) {
	query := `
		UPDATE issues
		SET title = COALESCE($2, title),
		    description = COALESCE($3, description),
		    status = COALESCE($4, status),
		    priority = COALESCE($5, priority),
		    type = COALESCE($6, type),
		    assignee_id = COALESCE($7, assignee_id),
		    sort_order = COALESCE($8, sort_order),
		    parent_issue_id = COALESCE($9, parent_issue_id),
		    due_date = COALESCE($10, due_date)
		WHERE id = $1
		RETURNING id, project_id, issue_key, title, description, status, priority, type,
		          assignee_id, reporter_id, sort_order, parent_issue_id, due_date,
		          created_at, updated_at
	`

	var issue models.Issue
	err := r.db.QueryRow(
		ctx,
		query,
		id,
		params.Title,
		params.Description,
		params.Status,
		params.Priority,
		params.Type,
		params.AssigneeID,
		params.SortOrder,
		params.ParentIssueID,
		params.DueDate,
	).Scan(
		&issue.ID, &issue.ProjectID, &issue.IssueKey, &issue.Title, &issue.Description,
		&issue.Status, &issue.Priority, &issue.Type, &issue.AssigneeID, &issue.ReporterID,
		&issue.SortOrder, &issue.ParentIssueID, &issue.DueDate, &issue.CreatedAt, &issue.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &issue, nil
}

// Delete hard deletes an issue
func (r *IssueRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM issues WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetSubIssues retrieves all child issues of a parent issue
func (r *IssueRepository) GetSubIssues(ctx context.Context, parentID uuid.UUID) ([]models.Issue, error) {
	query := `
		SELECT id, project_id, issue_key, title, description, status, priority, type,
		       assignee_id, reporter_id, sort_order, parent_issue_id, due_date,
		       created_at, updated_at
		FROM issues
		WHERE parent_issue_id = $1
		ORDER BY sort_order ASC, created_at DESC
	`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.Issue
	for rows.Next() {
		var issue models.Issue
		if err := rows.Scan(
			&issue.ID, &issue.ProjectID, &issue.IssueKey, &issue.Title, &issue.Description,
			&issue.Status, &issue.Priority, &issue.Type, &issue.AssigneeID, &issue.ReporterID,
			&issue.SortOrder, &issue.ParentIssueID, &issue.DueDate, &issue.CreatedAt, &issue.UpdatedAt,
		); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
