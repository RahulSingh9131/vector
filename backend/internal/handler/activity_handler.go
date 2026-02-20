package handler

import (
	"net/http"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ActivityHandler struct {
	Handler
	services *service.Services
}

func NewActivityHandler(s *server.Server, services *service.Services) *ActivityHandler {
	return &ActivityHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// ListIssueActivity returns the activity timeline for an issue
func (h *ActivityHandler) ListIssueActivity(c echo.Context, req *ListIssueActivityRequest) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	issueID, _ := uuid.Parse(req.IssueID)

	return h.services.Activity.ListByIssue(c.Request().Context(), issueID, req.Page, req.Limit)
}

// ListProjectActivity returns the activity feed for a project
func (h *ActivityHandler) ListProjectActivity(c echo.Context, req *ListProjectActivityRequest) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	return h.services.Activity.ListByProject(c.Request().Context(), projectID, req.Page, req.Limit)
}

// RegisterIssueRoutes registers issue activity routes
func (h *ActivityHandler) RegisterIssueRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.ListIssueActivity, http.StatusOK, &ListIssueActivityRequest{}))
}

// RegisterProjectRoutes registers project activity routes
func (h *ActivityHandler) RegisterProjectRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.ListProjectActivity, http.StatusOK, &ListProjectActivityRequest{}))
}
