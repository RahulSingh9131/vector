package handler

import (
	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	Handler
	services *service.Services
}

func NewOrganizationHandler(s *server.Server, services *service.Services) *OrganizationHandler {
	return &OrganizationHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// ListOrganizations lists all organizations the current user belongs to
func (h *OrganizationHandler) ListOrganizations(c echo.Context, req *EmptyRequest) ([]models.OrganizationMemberWithDetails, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	memberships, err := h.services.Organization.GetUserOrganizations(c.Request().Context(), userID)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

// GetOrganization returns details of a specific organization
func (h *OrganizationHandler) GetOrganization(c echo.Context, req *GetOrganizationRequest) (*models.Organization, error) {
	orgID, _ := uuid.Parse(req.ID) // already validated as UUID by request struct

	org, err := h.services.Organization.GetByID(c.Request().Context(), orgID)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// ListMembers lists all members of a specific organization
func (h *OrganizationHandler) ListMembers(c echo.Context, req *ListMembersRequest) ([]models.OrganizationMemberWithDetails, error) {
	orgID, _ := uuid.Parse(req.ID) // already validated as UUID by request struct

	members, err := h.services.Organization.GetMembers(c.Request().Context(), orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// CreateOrganization creates a new organization and adds the creator as owner
func (h *OrganizationHandler) CreateOrganization(c echo.Context, req *CreateOrganizationRequest) (*models.Organization, error) {
	creatorID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Organization.CreateOrganization(c.Request().Context(), creatorID, models.CreateOrganizationParams{
		Name:    req.Name,
		Slug:    req.Slug,
		LogoURL: req.LogoURL,
	})
}

// AddMember adds a user to an organization
func (h *OrganizationHandler) AddMember(c echo.Context, req *AddMemberRequest) (*models.OrganizationMember, error) {
	orgID, _ := uuid.Parse(req.ID)
	userID, _ := uuid.Parse(req.UserID)

	return h.services.Organization.AddMember(c.Request().Context(), orgID, userID, req.Role)
}

// UpdateMemberRole updates a member's role in an organization
func (h *OrganizationHandler) UpdateMemberRole(c echo.Context, req *UpdateMemberRoleRequest) (*models.OrganizationMember, error) {
	orgID, _ := uuid.Parse(req.ID)
	userID, _ := uuid.Parse(req.UserID)

	return h.services.Organization.UpdateMemberRole(c.Request().Context(), orgID, userID, req.Role)
}

// RemoveMember removes a user from an organization
func (h *OrganizationHandler) RemoveMember(c echo.Context, req *RemoveMemberRequest) (*EmptyResponse, error) {
	orgID, _ := uuid.Parse(req.ID)
	userID, _ := uuid.Parse(req.UserID)

	if err := h.services.Organization.RemoveMember(c.Request().Context(), orgID, userID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}
