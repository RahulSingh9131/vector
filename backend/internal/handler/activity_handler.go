package handler

import (
	"net/http"
	"time"

	"github.com/RahulSingh9131/vector/internal/errs"
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

// parseActivityFilters builds ActivityFilters from the common filter query params.
func parseActivityFilters(action, entityType, actorID, from, to *string) (models.ActivityFilters, error) {
	var filters models.ActivityFilters

	filters.Action = action
	filters.EntityType = entityType

	if actorID != nil {
		parsed, err := uuid.Parse(*actorID)
		if err != nil {
			return filters, errs.NewBadRequestError("invalid actor_id: must be a valid UUID", true, nil, nil, nil)
		}
		filters.ActorID = &parsed
	}

	if from != nil {
		t, err := time.Parse(time.RFC3339, *from)
		if err != nil {
			return filters, errs.NewBadRequestError("invalid 'from': must be RFC3339 format (e.g. 2026-02-01T00:00:00Z)", true, nil, nil, nil)
		}
		filters.From = &t
	}

	if to != nil {
		t, err := time.Parse(time.RFC3339, *to)
		if err != nil {
			return filters, errs.NewBadRequestError("invalid 'to': must be RFC3339 format (e.g. 2026-02-28T23:59:59Z)", true, nil, nil, nil)
		}
		filters.To = &t
	}

	return filters, nil
}

// ListIssueActivity returns the activity timeline for an issue
func (h *ActivityHandler) ListIssueActivity(c echo.Context, req *ListIssueActivityRequest) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	issueID, _ := uuid.Parse(req.IssueID)

	filters, err := parseActivityFilters(req.Action, req.EntityType, req.ActorID, req.From, req.To)
	if err != nil {
		return nil, err
	}

	return h.services.Activity.ListByIssue(c.Request().Context(), issueID, req.Page, req.Limit, filters)
}

// ListProjectActivity returns the activity feed for a project
func (h *ActivityHandler) ListProjectActivity(c echo.Context, req *ListProjectActivityRequest) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	filters, err := parseActivityFilters(req.Action, req.EntityType, req.ActorID, req.From, req.To)
	if err != nil {
		return nil, err
	}

	return h.services.Activity.ListByProject(c.Request().Context(), projectID, req.Page, req.Limit, filters)
}

// RegisterIssueRoutes registers issue activity routes
func (h *ActivityHandler) RegisterIssueRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.ListIssueActivity, http.StatusOK, &ListIssueActivityRequest{}))
}

// RegisterProjectRoutes registers project activity routes
func (h *ActivityHandler) RegisterProjectRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.ListProjectActivity, http.StatusOK, &ListProjectActivityRequest{}))
}

