package service

import (
	"context"
	"strings"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	"github.com/google/uuid"
)

// Allowed project member roles
var allowedProjectRoles = map[string]bool{
	"admin":  true,
	"member": true,
	"viewer": true,
}

// Allowed project statuses
var allowedProjectStatuses = map[string]bool{
	"active":   true,
	"archived": true,
}

// ProjectService handles business logic for projects
type ProjectService struct {
	server        *server.Server
	projectRepo   *repository.ProjectRepository
	memberRepo    *repository.ProjectMemberRepository
	orgRepo       *repository.OrganizationRepository
	orgMemberRepo *repository.OrganizationMemberRepository
}

// NewProjectService creates a new project service
func NewProjectService(s *server.Server, repos *repository.Repositories) *ProjectService {
	return &ProjectService{
		server:        s,
		projectRepo:   repos.Project,
		memberRepo:    repos.ProjectMember,
		orgRepo:       repos.Organization,
		orgMemberRepo: repos.OrganizationMember,
	}
}

// CreateProject creates a new project within an organization
func (s *ProjectService) CreateProject(ctx context.Context, orgID, userID uuid.UUID, params models.CreateProjectParams) (*models.Project, error) {
	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("name", params.Name).
		Msg("creating project")

	// Check org exists
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if org == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Check project limit
	projectCount, err := s.projectRepo.GetProjectCount(ctx, orgID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	if projectCount >= org.MaxProjects {
		return nil, errs.NewBadRequestError(
			"Organization has reached the maximum project limit",
			true, nil, nil, nil,
		)
	}

	// Validate identifier (uppercase, alphanumeric, 2-10 chars)
	params.Identifier = strings.ToUpper(params.Identifier)

	// Set org and creator
	params.OrganizationID = orgID
	params.CreatedBy = userID

	project, err := s.projectRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("name", params.Name).
			Msg("failed to create project")
		return nil, sqlerr.HandleError(err)
	}

	// Auto-add the creator as project admin
	_, err = s.memberRepo.AddMember(ctx, models.CreateProjectMemberParams{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      "admin",
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", project.ID.String()).
			Str("user_id", userID.String()).
			Msg("failed to add creator as project admin")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", project.ID.String()).
		Str("org_id", orgID.String()).
		Str("name", params.Name).
		Msg("project created successfully")

	return project, nil
}

// GetByID retrieves a project by ID
func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if project == nil {
		return nil, errs.NewNotFoundError("Project not found", false, nil)
	}

	return project, nil
}

// ListProjects retrieves projects the user is a member of within an organization
func (s *ProjectService) ListProjects(ctx context.Context, orgID, userID uuid.UUID) ([]models.Project, error) {
	projects, err := s.projectRepo.ListByUserAndOrganization(ctx, userID, orgID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to list user projects")
		return nil, sqlerr.HandleError(err)
	}

	return projects, nil
}

// UpdateProject updates a project's information
func (s *ProjectService) UpdateProject(ctx context.Context, projectID uuid.UUID, params models.UpdateProjectParams) (*models.Project, error) {
	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Msg("updating project")

	// Validate status if provided
	if params.Status != nil && !allowedProjectStatuses[*params.Status] {
		return nil, errs.NewBadRequestError(
			"Invalid project status. Must be one of: active, archived",
			true, nil,
			[]errs.FieldError{{Field: "status", Error: "must be one of: active, archived"}},
			nil,
		)
	}

	project, err := s.projectRepo.Update(ctx, projectID, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Msg("failed to update project")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", projectID.String()).
		Msg("project updated successfully")

	return project, nil
}

// DeleteProject soft-deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Msg("deleting project")

	if err := s.projectRepo.Delete(ctx, projectID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Msg("failed to delete project")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", projectID.String()).
		Msg("project deleted successfully")

	return nil
}

// AddMember adds a user to a project
func (s *ProjectService) AddMember(ctx context.Context, projectID, userID uuid.UUID, role string) (*models.ProjectMember, error) {
	// Validate role
	if !allowedProjectRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: admin, member, viewer",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: admin, member, viewer"}},
			nil,
		)
	}

	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("adding member to project")

	// Verify the project exists
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if project == nil {
		return nil, errs.NewNotFoundError("Project not found", false, nil)
	}

	// Verify the user is a member of the organization
	orgMember, err := s.orgMemberRepo.GetMember(ctx, project.OrganizationID, userID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if orgMember == nil {
		return nil, errs.NewBadRequestError(
			"User is not a member of this organization",
			true, nil, nil, nil,
		)
	}

	member, err := s.memberRepo.AddMember(ctx, models.CreateProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Str("user_id", userID.String()).
			Msg("failed to add member to project")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("member added to project successfully")

	return member, nil
}

// RemoveMember removes a user from a project
func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Msg("removing member from project")

	// Check if this would remove the last admin
	member, err := s.memberRepo.GetMember(ctx, projectID, userID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if member == nil {
		return errs.NewNotFoundError("Member not found in project", false, nil)
	}

	if member.Role == "admin" {
		adminCount, err := s.memberRepo.GetAdminCount(ctx, projectID)
		if err != nil {
			return sqlerr.HandleError(err)
		}
		if adminCount <= 1 {
			return errs.NewBadRequestError(
				"Cannot remove the last admin from the project",
				true, nil, nil, nil,
			)
		}
	}

	if err := s.memberRepo.RemoveMember(ctx, projectID, userID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Str("user_id", userID.String()).
			Msg("failed to remove member from project")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Msg("member removed from project successfully")

	return nil
}

// UpdateMemberRole updates a member's role in a project
func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectID, userID uuid.UUID, role string) (*models.ProjectMember, error) {
	// Validate role
	if !allowedProjectRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: admin, member, viewer",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: admin, member, viewer"}},
			nil,
		)
	}

	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("updating member role")

	// If demoting from admin, check this wouldn't remove the last admin
	existingMember, err := s.memberRepo.GetMember(ctx, projectID, userID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingMember == nil {
		return nil, errs.NewNotFoundError("Member not found in project", false, nil)
	}

	if existingMember.Role == "admin" && role != "admin" {
		adminCount, adminErr := s.memberRepo.GetAdminCount(ctx, projectID)
		if adminErr != nil {
			return nil, sqlerr.HandleError(adminErr)
		}
		if adminCount <= 1 {
			return nil, errs.NewBadRequestError(
				"Cannot demote the last admin of the project",
				true, nil, nil, nil,
			)
		}
	}

	member, err := s.memberRepo.UpdateRole(ctx, projectID, userID, role)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Str("user_id", userID.String()).
			Msg("failed to update member role")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("project_id", projectID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("member role updated successfully")

	return member, nil
}

// GetMembers retrieves all members of a project
func (s *ProjectService) GetMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMemberWithDetails, error) {
	members, err := s.memberRepo.GetMembersByProject(ctx, projectID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Msg("failed to get project members")
		return nil, sqlerr.HandleError(err)
	}

	return members, nil
}

// CheckMembership checks if a user is a member of a project
func (s *ProjectService) CheckMembership(ctx context.Context, projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	member, err := s.memberRepo.GetMember(ctx, projectID, userID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	return member, nil
}
