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

type LabelHandler struct {
	Handler
	services *service.Services
}

func NewLabelHandler(s *server.Server, services *service.Services) *LabelHandler {
	return &LabelHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// CreateLabel creates a new label within a project
func (h *LabelHandler) CreateLabel(c echo.Context, req *CreateLabelRequest) (*models.Label, error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	return h.services.Label.CreateLabel(c.Request().Context(), projectID, models.CreateLabelParams{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	})
}

// ListLabels lists all labels for a project
func (h *LabelHandler) ListLabels(c echo.Context, req *ListLabelsRequest) ([]models.Label, error) {
	projectID, _ := uuid.Parse(req.ProjectID)

	return h.services.Label.ListByProject(c.Request().Context(), projectID)
}

// GetLabel returns a specific label
func (h *LabelHandler) GetLabel(c echo.Context, req *GetLabelRequest) (*models.Label, error) {
	labelID, _ := uuid.Parse(req.LabelID)

	return h.services.Label.GetByID(c.Request().Context(), labelID)
}

// UpdateLabel updates a label
func (h *LabelHandler) UpdateLabel(c echo.Context, req *UpdateLabelRequest) (*models.Label, error) {
	labelID, _ := uuid.Parse(req.LabelID)

	return h.services.Label.UpdateLabel(c.Request().Context(), labelID, models.UpdateLabelParams{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	})
}

// DeleteLabel deletes a label
func (h *LabelHandler) DeleteLabel(c echo.Context, req *DeleteLabelRequest) (*EmptyResponse, error) {
	labelID, _ := uuid.Parse(req.LabelID)

	if err := h.services.Label.DeleteLabel(c.Request().Context(), labelID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// GetIssueLabels returns all labels attached to an issue
func (h *LabelHandler) GetIssueLabels(c echo.Context, req *GetIssueLabelsRequest) ([]models.Label, error) {
	issueID, _ := uuid.Parse(req.IssueID)

	return h.services.Label.GetLabelsByIssue(c.Request().Context(), issueID)
}

// AddLabelToIssue attaches a label to an issue
func (h *LabelHandler) AddLabelToIssue(c echo.Context, req *AddLabelToIssueRequest) (*EmptyResponse, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	labelID, _ := uuid.Parse(req.LabelID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Label.AddLabelToIssue(c.Request().Context(), issueID, labelID, userID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// RemoveLabelFromIssue removes a label from an issue
func (h *LabelHandler) RemoveLabelFromIssue(c echo.Context, req *RemoveLabelFromIssueRequest) (*EmptyResponse, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	labelID, _ := uuid.Parse(req.LabelID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Label.RemoveLabelFromIssue(c.Request().Context(), issueID, labelID, userID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// RegisterRoutes registers label CRUD routes
func (h *LabelHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", Handle(h.Handler, h.CreateLabel, http.StatusCreated, &CreateLabelRequest{}))
	g.GET("", Handle(h.Handler, h.ListLabels, http.StatusOK, &ListLabelsRequest{}))
	g.GET("/:labelId", Handle(h.Handler, h.GetLabel, http.StatusOK, &GetLabelRequest{}))
	g.PATCH("/:labelId", Handle(h.Handler, h.UpdateLabel, http.StatusOK, &UpdateLabelRequest{}))
	g.DELETE("/:labelId", Handle(h.Handler, h.DeleteLabel, http.StatusNoContent, &DeleteLabelRequest{}))
}

// RegisterIssueLabelRoutes registers issue-label association routes
func (h *LabelHandler) RegisterIssueLabelRoutes(g *echo.Group) {
	g.GET("", Handle(h.Handler, h.GetIssueLabels, http.StatusOK, &GetIssueLabelsRequest{}))
	g.POST("", Handle(h.Handler, h.AddLabelToIssue, http.StatusCreated, &AddLabelToIssueRequest{}))
	g.DELETE("/:labelId", Handle(h.Handler, h.RemoveLabelFromIssue, http.StatusNoContent, &RemoveLabelFromIssueRequest{}))
}
