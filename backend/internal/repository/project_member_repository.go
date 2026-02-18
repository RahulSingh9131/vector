package repository

import (
	"context"
	"errors"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ProjectMemberRepository handles database operations for project members
type ProjectMemberRepository struct {
	server *server.Server
}

// NewProjectMemberRepository creates a new project member repository
func NewProjectMemberRepository(s *server.Server) *ProjectMemberRepository {
	return &ProjectMemberRepository{
		server: s,
	}
}

// AddMember adds a user to a project
func (r *ProjectMemberRepository) AddMember(ctx context.Context, params models.CreateProjectMemberParams) (*models.ProjectMember, error) {
	query := `
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, project_id, user_id, role, joined_at
	`

	var member models.ProjectMember
	err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.ProjectID,
		params.UserID,
		params.Role,
	).Scan(
		&member.ID,
		&member.ProjectID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

// RemoveMember removes a user from a project
func (r *ProjectMemberRepository) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	query := `
		DELETE FROM project_members
		WHERE project_id = $1 AND user_id = $2
	`

	_, err := r.server.DB.Pool.Exec(ctx, query, projectID, userID)
	return err
}

// UpdateRole updates a member's role in a project
func (r *ProjectMemberRepository) UpdateRole(ctx context.Context, projectID, userID uuid.UUID, role string) (*models.ProjectMember, error) {
	query := `
		UPDATE project_members
		SET role = $3
		WHERE project_id = $1 AND user_id = $2
		RETURNING id, project_id, user_id, role, joined_at
	`

	var member models.ProjectMember
	err := r.server.DB.Pool.QueryRow(ctx, query, projectID, userID, role).Scan(
		&member.ID,
		&member.ProjectID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

// GetMember retrieves a specific member of a project
func (r *ProjectMemberRepository) GetMember(ctx context.Context, projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	query := `
		SELECT id, project_id, user_id, role, joined_at
		FROM project_members
		WHERE project_id = $1 AND user_id = $2
	`

	var member models.ProjectMember
	err := r.server.DB.Pool.QueryRow(ctx, query, projectID, userID).Scan(
		&member.ID,
		&member.ProjectID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}

// GetMembersByProject retrieves all members of a project with user details
func (r *ProjectMemberRepository) GetMembersByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMemberWithDetails, error) {
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.joined_at,
		       u.id, u.clerk_user_id, u.email, u.first_name, u.last_name, u.avatar_url,
		       u.is_active, u.last_login_at, u.created_at, u.updated_at
		FROM project_members pm
		INNER JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.joined_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMemberWithDetails
	for rows.Next() {
		var m models.ProjectMemberWithDetails
		var user models.User
		if err := rows.Scan(
			&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.JoinedAt,
			&user.ID, &user.ClerkUserID, &user.Email, &user.FirstName, &user.LastName,
			&user.AvatarURL, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.User = &user
		members = append(members, m)
	}

	return members, nil
}

// GetProjectsByUser retrieves all projects a user is a member of within an organization
func (r *ProjectMemberRepository) GetProjectsByUser(ctx context.Context, userID, orgID uuid.UUID) ([]models.ProjectMemberWithDetails, error) {
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.joined_at,
		       p.id, p.organization_id, p.name, p.slug, p.description, p.status,
		       p.identifier, p.issue_counter, p.created_by, p.created_at, p.updated_at
		FROM project_members pm
		INNER JOIN projects p ON p.id = pm.project_id
		WHERE pm.user_id = $1 AND p.organization_id = $2 AND p.status != 'deleted'
		ORDER BY p.created_at DESC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, userID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ProjectMemberWithDetails
	for rows.Next() {
		var m models.ProjectMemberWithDetails
		var project models.Project
		if err := rows.Scan(
			&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.JoinedAt,
			&project.ID, &project.OrganizationID, &project.Name, &project.Slug,
			&project.Description, &project.Status, &project.Identifier,
			&project.IssueCounter, &project.CreatedBy, &project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		m.Project = &project
		members = append(members, m)
	}

	return members, nil
}

// GetMemberCount returns the number of members in a project
func (r *ProjectMemberRepository) GetMemberCount(ctx context.Context, projectID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM project_members
		WHERE project_id = $1
	`

	var count int
	err := r.server.DB.Pool.QueryRow(ctx, query, projectID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetAdminCount returns the number of admins in a project
func (r *ProjectMemberRepository) GetAdminCount(ctx context.Context, projectID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM project_members
		WHERE project_id = $1 AND role = 'admin'
	`

	var count int
	err := r.server.DB.Pool.QueryRow(ctx, query, projectID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
