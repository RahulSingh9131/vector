package handler

import (
	"fmt"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// SearchHandler handles search-related requests
type SearchHandler struct {
	Handler
	server   *server.Server
	services *service.Services
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(s *server.Server, services *service.Services) *SearchHandler {
	return &SearchHandler{
		Handler:  NewHandler(s),
		server:   s,
		services: services,
	}
}

// SearchRequest defines the request parameters for search
type SearchRequest struct {
	Query     string  `query:"query"`
	Type      *string `query:"type"`
	ProjectID *string `query:"projectId"`
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
}

// Validate validates the search request
func (r *SearchRequest) Validate() error {
	if r.Query == "" {
		return fmt.Errorf("query is required")
	}
	if r.ProjectID != nil && *r.ProjectID != "" {
		if _, err := uuid.Parse(*r.ProjectID); err != nil {
			return fmt.Errorf("invalid projectId")
		}
	}
	return nil
}

// Search performs a unified search across issues and comments
func (h *SearchHandler) Search(c echo.Context, req *SearchRequest) ([]models.SearchResult, error) {
	params := models.SearchParams{
		Query:  req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if req.Type != nil && *req.Type != "" {
		t := models.SearchResultType(*req.Type)
		params.Type = &t
	}

	if req.ProjectID != nil && *req.ProjectID != "" {
		id, _ := uuid.Parse(*req.ProjectID)
		params.ProjectID = &id
	}

	return h.services.Search.Search(c.Request().Context(), params)
}
