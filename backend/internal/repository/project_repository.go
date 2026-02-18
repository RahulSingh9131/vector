package repository

import (
	"context"
	"errors"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	server *server.Server
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(s *server.Server) *ProjectRepository {
	return &ProjectRepository{
		server: s,
	}
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, params models.CreateProjectParams) (*models.Project, error) {
	query := `
		INSERT INTO projects (organization_id, name, slug, description, identifier, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, organization_id, name, slug, description, status, identifier,
		          issue_counter, created_by, created_at, updated_at
	`

	var project models.Project
	err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.OrganizationID,
		params.Name,
		params.Slug,
		params.Description,
		params.Identifier,
		params.CreatedBy,
	).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.Identifier,
		&project.IssueCounter,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetByID retrieves a project by its UUID
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
		SELECT id, organization_id, name, slug, description, status, identifier,
		       issue_counter, created_by, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	var project models.Project
	err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.Identifier,
		&project.IssueCounter,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

// GetBySlug retrieves a project by organization ID and slug
func (r *ProjectRepository) GetBySlug(ctx context.Context, orgID uuid.UUID, slug string) (*models.Project, error) {
	query := `
		SELECT id, organization_id, name, slug, description, status, identifier,
		       issue_counter, created_by, created_at, updated_at
		FROM projects
		WHERE organization_id = $1 AND slug = $2
	`

	var project models.Project
	err := r.server.DB.Pool.QueryRow(ctx, query, orgID, slug).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.Identifier,
		&project.IssueCounter,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &project, nil
}

// ListByOrganization retrieves all projects for an organization
func (r *ProjectRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	query := `
		SELECT id, organization_id, name, slug, description, status, identifier,
		       issue_counter, created_by, created_at, updated_at
		FROM projects
		WHERE organization_id = $1 AND status != 'deleted'
		ORDER BY created_at DESC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID,
			&project.OrganizationID,
			&project.Name,
			&project.Slug,
			&project.Description,
			&project.Status,
			&project.Identifier,
			&project.IssueCounter,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// ListByUserAndOrganization retrieves all projects the user is a member of within an organization
func (r *ProjectRepository) ListByUserAndOrganization(ctx context.Context, userID, orgID uuid.UUID) ([]models.Project, error) {
	query := `
		SELECT p.id, p.organization_id, p.name, p.slug, p.description, p.status, p.identifier,
		       p.issue_counter, p.created_by, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN project_members pm ON pm.project_id = p.id
		WHERE p.organization_id = $1 AND pm.user_id = $2 AND p.status != 'deleted'
		ORDER BY p.created_at DESC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, orgID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID,
			&project.OrganizationID,
			&project.Name,
			&project.Slug,
			&project.Description,
			&project.Status,
			&project.Identifier,
			&project.IssueCounter,
			&project.CreatedBy,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateProjectParams) (*models.Project, error) {
	query := `
		UPDATE projects
		SET name = COALESCE($2, name),
		    slug = COALESCE($3, slug),
		    description = COALESCE($4, description),
		    status = COALESCE($5, status)
		WHERE id = $1
		RETURNING id, organization_id, name, slug, description, status, identifier,
		          issue_counter, created_by, created_at, updated_at
	`

	var project models.Project
	err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		id,
		params.Name,
		params.Slug,
		params.Description,
		params.Status,
	).Scan(
		&project.ID,
		&project.OrganizationID,
		&project.Name,
		&project.Slug,
		&project.Description,
		&project.Status,
		&project.Identifier,
		&project.IssueCounter,
		&project.CreatedBy,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// Delete soft deletes a project by setting status to 'deleted'
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE projects
		SET status = 'deleted'
		WHERE id = $1
	`

	_, err := r.server.DB.Pool.Exec(ctx, query, id)
	return err
}

// IncrementIssueCounter atomically increments and returns the new issue counter
func (r *ProjectRepository) IncrementIssueCounter(ctx context.Context, id uuid.UUID) (int, error) {
	query := `
		UPDATE projects
		SET issue_counter = issue_counter + 1
		WHERE id = $1
		RETURNING issue_counter
	`

	var counter int
	err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(&counter)
	if err != nil {
		return 0, err
	}

	return counter, nil
}

// GetProjectCount returns the number of active projects in an organization
func (r *ProjectRepository) GetProjectCount(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM projects
		WHERE organization_id = $1 AND status != 'deleted'
	`

	var count int
	err := r.server.DB.Pool.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
