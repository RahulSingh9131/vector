package handler

import (
	"time"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	Handler
	services *service.Services
}

func NewProjectHandler(s *server.Server, services *service.Services) *ProjectHandler {
	return &ProjectHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// CreateProject creates a new project within an organization
func (h *ProjectHandler) CreateProject(c echo.Context, req *CreateProjectRequest) (*models.Project, error) {
	orgID, _ := uuid.Parse(req.OrgID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Project.CreateProject(c.Request().Context(), orgID, userID, models.CreateProjectParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Identifier:  req.Identifier,
	})
}

// ListProjects lists all projects the current user is a member of in an organization
func (h *ProjectHandler) ListProjects(c echo.Context, req *ListProjectsRequest) ([]models.Project, error) {
	orgID, _ := uuid.Parse(req.OrgID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Project.ListProjects(c.Request().Context(), orgID, userID)
}

// GetProject returns details of a specific project
func (h *ProjectHandler) GetProject(c echo.Context, req *GetProjectRequest) (*models.Project, error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	return h.services.Project.GetByID(c.Request().Context(), projectID)
}

// UpdateProject updates a project
func (h *ProjectHandler) UpdateProject(c echo.Context, req *UpdateProjectRequest) (*models.Project, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	actorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Project.UpdateProject(c.Request().Context(), projectID, actorID, models.UpdateProjectParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      req.Status,
	})
}

// DeleteProject soft-deletes a project
func (h *ProjectHandler) DeleteProject(c echo.Context, req *DeleteProjectRequest) (*EmptyResponse, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	actorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Project.DeleteProject(c.Request().Context(), projectID, actorID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// ListMembers lists all members of a project
func (h *ProjectHandler) ListMembers(c echo.Context, req *ListProjectMembersRequest) ([]models.ProjectMemberWithDetails, error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	return h.services.Project.GetMembers(c.Request().Context(), projectID)
}

// AddMember adds a user to a project
func (h *ProjectHandler) AddMember(c echo.Context, req *AddProjectMemberRequest) (*models.ProjectMember, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	userID, _ := uuid.Parse(req.UserID)
	actorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Project.AddMember(c.Request().Context(), projectID, userID, actorID, req.Role)
}

// UpdateMemberRole updates a member's role in a project
func (h *ProjectHandler) UpdateMemberRole(c echo.Context, req *UpdateProjectMemberRoleRequest) (*models.ProjectMember, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	userID, _ := uuid.Parse(req.UserID)
	actorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Project.UpdateMemberRole(c.Request().Context(), projectID, userID, actorID, req.Role)
}

// RemoveMember removes a user from a project
func (h *ProjectHandler) RemoveMember(c echo.Context, req *RemoveProjectMemberRequest) (*EmptyResponse, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	userID, _ := uuid.Parse(req.UserID)
	actorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Project.RemoveMember(c.Request().Context(), projectID, userID, actorID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}


// parseTime parses a time string in RFC3339 format
func parseTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
