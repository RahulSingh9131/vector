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

type IssueHandler struct {
	Handler
	services *service.Services
}

func NewIssueHandler(s *server.Server, services *service.Services) *IssueHandler {
	return &IssueHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// CreateIssue creates a new issue within a project
func (h *IssueHandler) CreateIssue(c echo.Context, req *CreateIssueRequest) (*models.Issue, error) {
	projectID, _ := uuid.Parse(req.ProjectID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	params := models.CreateIssueParams{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Type:        req.Type,
	}

	if req.AssigneeID != nil {
		id, _ := uuid.Parse(*req.AssigneeID)
		params.AssigneeID = &id
	}

	if req.ParentIssueID != nil {
		id, _ := uuid.Parse(*req.ParentIssueID)
		params.ParentIssueID = &id
	}

	if req.DueDate != nil {
		t := parseTime(*req.DueDate)
		params.DueDate = t
	}

	return h.services.Issue.CreateIssue(c.Request().Context(), projectID, userID, params)
}

// ListIssues lists issues for a project with optional filters
func (h *IssueHandler) ListIssues(c echo.Context, req *ListIssuesRequest) (*models.PaginatedResponse[models.IssueWithDetails], error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	filters := models.IssueFilters{
		Status:   req.Status,
		Priority: req.Priority,
		Type:     req.Type,
		Page:     req.Page,
		Limit:    req.Limit,
	}

	if req.AssigneeID != nil {
		id, _ := uuid.Parse(*req.AssigneeID)
		filters.AssigneeID = &id
	}

	return h.services.Issue.ListIssues(c.Request().Context(), projectID, filters)
}

// GetIssue returns details of a specific issue
func (h *IssueHandler) GetIssue(c echo.Context, req *GetIssueRequest) (*models.IssueWithDetails, error) {
	issueID, _ := uuid.Parse(req.IssueID)

	return h.services.Issue.GetByID(c.Request().Context(), issueID)
}

// UpdateIssue updates an issue
func (h *IssueHandler) UpdateIssue(c echo.Context, req *UpdateIssueRequest) (*models.Issue, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	params := models.UpdateIssueParams{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Type:        req.Type,
		SortOrder:   req.SortOrder,
	}

	if req.AssigneeID != nil {
		id, _ := uuid.Parse(*req.AssigneeID)
		params.AssigneeID = &id
	}

	if req.ParentIssueID != nil {
		id, _ := uuid.Parse(*req.ParentIssueID)
		params.ParentIssueID = &id
	}

	if req.DueDate != nil {
		t := parseTime(*req.DueDate)
		params.DueDate = t
	}

	return h.services.Issue.UpdateIssue(c.Request().Context(), issueID, userID, params)
}

// DeleteIssue deletes an issue
func (h *IssueHandler) DeleteIssue(c echo.Context, req *DeleteIssueRequest) (*EmptyResponse, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Issue.DeleteIssue(c.Request().Context(), issueID, userID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// AssignIssue assigns or unassigns an issue
func (h *IssueHandler) AssignIssue(c echo.Context, req *AssignIssueRequest) (*models.Issue, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != nil {
		id, _ := uuid.Parse(*req.AssigneeID)
		assigneeID = &id
	}

	return h.services.Issue.AssignIssue(c.Request().Context(), issueID, userID, assigneeID)
}

// UpdateIssueStatus updates an issue's status
func (h *IssueHandler) UpdateIssueStatus(c echo.Context, req *UpdateIssueStatusRequest) (*models.Issue, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Issue.UpdateStatus(c.Request().Context(), issueID, userID, req.Status)
}

// RegisterRoutes registers all issue routes
func (h *IssueHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", Handle(h.Handler, h.CreateIssue, http.StatusCreated, &CreateIssueRequest{}))
	g.GET("", Handle(h.Handler, h.ListIssues, http.StatusOK, &ListIssuesRequest{}))
	g.GET("/:issueId", Handle(h.Handler, h.GetIssue, http.StatusOK, &GetIssueRequest{}))
	g.PATCH("/:issueId", Handle(h.Handler, h.UpdateIssue, http.StatusOK, &UpdateIssueRequest{}))
	g.DELETE("/:issueId", Handle(h.Handler, h.DeleteIssue, http.StatusNoContent, &DeleteIssueRequest{}))
	g.PATCH("/:issueId/assign", Handle(h.Handler, h.AssignIssue, http.StatusOK, &AssignIssueRequest{}))
	g.PATCH("/:issueId/status", Handle(h.Handler, h.UpdateIssueStatus, http.StatusOK, &UpdateIssueStatusRequest{}))
}
