package service

import (
	"context"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
)

// SearchService handles search-related logic
type SearchService struct {
	server     *server.Server
	searchRepo *repository.SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(s *server.Server, repos *repository.Repositories) *SearchService {
	return &SearchService{
		server:     s,
		searchRepo: repos.Search,
	}
}

// Search performs a unified search across issues and comments
func (s *SearchService) Search(ctx context.Context, params models.SearchParams) ([]models.SearchResult, error) {
	if params.Limit == 0 {
		params.Limit = 20
	}
	return s.searchRepo.Search(ctx, params)
}
