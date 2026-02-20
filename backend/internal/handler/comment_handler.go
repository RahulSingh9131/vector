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

type CommentHandler struct {
	Handler
	services *service.Services
}

func NewCommentHandler(s *server.Server, services *service.Services) *CommentHandler {
	return &CommentHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

// CreateComment creates a new comment on an issue
func (h *CommentHandler) CreateComment(c echo.Context, req *CreateCommentRequest) (*models.CommentWithAuthor, error) {
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	params := models.CreateCommentParams{
		Body: req.Body,
	}

	if req.ParentCommentID != nil {
		id, _ := uuid.Parse(*req.ParentCommentID)
		params.ParentCommentID = &id
	}

	return h.services.Comment.CreateComment(c.Request().Context(), issueID, userID, params)
}

// ListComments lists threaded comments for an issue
func (h *CommentHandler) ListComments(c echo.Context, req *ListCommentsRequest) (*models.PaginatedResponse[models.CommentThread], error) {
	issueID, _ := uuid.Parse(req.IssueID)

	return h.services.Comment.ListByIssue(c.Request().Context(), issueID, req.Page, req.Limit)
}

// GetComment returns a single comment
func (h *CommentHandler) GetComment(c echo.Context, req *GetCommentRequest) (*models.CommentWithAuthor, error) {
	commentID, _ := uuid.Parse(req.CommentID)

	return h.services.Comment.GetByID(c.Request().Context(), commentID)
}

// UpdateComment edits a comment (author only)
func (h *CommentHandler) UpdateComment(c echo.Context, req *UpdateCommentRequest) (*models.CommentWithAuthor, error) {
	commentID, _ := uuid.Parse(req.CommentID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return h.services.Comment.UpdateComment(c.Request().Context(), commentID, userID, models.UpdateCommentParams{
		Body: req.Body,
	})
}

// DeleteComment soft-deletes a comment (author or project admin)
func (h *CommentHandler) DeleteComment(c echo.Context, req *DeleteCommentRequest) (*EmptyResponse, error) {
	commentID, _ := uuid.Parse(req.CommentID)
	projectID, _ := uuid.Parse(req.ProjectID)
	issueID, _ := uuid.Parse(req.IssueID)
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	if err := h.services.Comment.DeleteComment(c.Request().Context(), commentID, userID, projectID, issueID); err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

// RegisterRoutes registers comment routes
func (h *CommentHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", Handle(h.Handler, h.CreateComment, http.StatusCreated, &CreateCommentRequest{}))
	g.GET("", Handle(h.Handler, h.ListComments, http.StatusOK, &ListCommentsRequest{}))
	g.GET("/:commentId", Handle(h.Handler, h.GetComment, http.StatusOK, &GetCommentRequest{}))
	g.PATCH("/:commentId", Handle(h.Handler, h.UpdateComment, http.StatusOK, &UpdateCommentRequest{}))
	g.DELETE("/:commentId", Handle(h.Handler, h.DeleteComment, http.StatusNoContent, &DeleteCommentRequest{}))
}
