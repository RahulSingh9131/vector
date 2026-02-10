package handler

import (
	"net/http"

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

// RegisterRoutes registers all organization routes
func (h *OrganizationHandler) RegisterRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.ListOrganizations, http.StatusOK, &EmptyRequest{}))
	g.GET("/:id", Handle(h.Handler, h.GetOrganization, http.StatusOK, &GetOrganizationRequest{}))
	g.GET("/:id/members", Handle(h.Handler, h.ListMembers, http.StatusOK, &ListMembersRequest{}))
}
