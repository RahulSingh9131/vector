package handler

import (
	"net/http"

	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	server   *server.Server
	services *service.Services
}

func NewOrganizationHandler(s *server.Server, services *service.Services) *OrganizationHandler {
	return &OrganizationHandler{
		server:   s,
		services: services,
	}
}

// ListOrganizations lists all organizations the current user belongs to
func (h *OrganizationHandler) ListOrganizations(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "User ID not found in context")
	}

	memberships, err := h.services.Organization.GetUserOrganizations(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user organizations")
	}

	return c.JSON(http.StatusOK, memberships)
}

// GetOrganization returns details of a specific organization
func (h *OrganizationHandler) GetOrganization(c echo.Context) error {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	org, err := h.services.Organization.GetByID(c.Request().Context(), orgID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return c.JSON(http.StatusOK, org)
}

// ListMembers lists all members of a specific organization
func (h *OrganizationHandler) ListMembers(c echo.Context) error {
	orgIDStr := c.Param("id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	members, err := h.services.Organization.GetMembers(c.Request().Context(), orgID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch organization members")
	}

	return c.JSON(http.StatusOK, members)
}
