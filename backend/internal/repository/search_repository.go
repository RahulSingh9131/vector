package repository

import (
	"context"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
)

// SearchRepository handles unified full-text search across various entities
type SearchRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(s *server.Server) *SearchRepository {
	return &SearchRepository{
		db:     s.DB.Pool,
		server: s,
	}
}

// Search performs a unified search across issues and comments
func (r *SearchRepository) Search(ctx context.Context, params models.SearchParams) ([]models.SearchResult, error) {
	// We use websearch_to_tsquery for natural language search strings
	// and ts_rank_cd for ranking based on cover density
	
	issueQuery := `
		SELECT 
			id, 
			'issue' as type, 
			title, 
			COALESCE(description, '') as description, 
			project_id, 
			NULL::uuid as issue_id, 
			issue_key, 
			ts_rank_cd(search_vector, websearch_to_tsquery('english', $1)) as rank
		FROM issues
		WHERE search_vector @@ websearch_to_tsquery('english', $1)
	`

	commentQuery := `
		SELECT 
			c.id, 
			'comment' as type, 
			'' as title, 
			c.body as description, 
			i.project_id, 
			c.issue_id, 
			'' as issue_key, 
			ts_rank_cd(c.search_vector, websearch_to_tsquery('english', $1)) as rank
		FROM comments c
		JOIN issues i ON c.issue_id = i.id
		WHERE c.search_vector @@ websearch_to_tsquery('english', $1)
	`

	var finalQuery string
	args := []interface{}{params.Query}
	argIdx := 2

	// Filter by type if provided
	if params.Type != nil {
		if *params.Type == models.SearchResultTypeIssue {
			finalQuery = issueQuery
		} else if *params.Type == models.SearchResultTypeComment {
			finalQuery = commentQuery
		}
	}

	if finalQuery == "" {
		finalQuery = fmt.Sprintf("(%s) UNION ALL (%s)", issueQuery, commentQuery)
	}

	// Filter by ProjectID if provided
	if params.ProjectID != nil {
		finalQuery = fmt.Sprintf("SELECT * FROM (%s) AS results WHERE project_id = $%d", finalQuery, argIdx)
		args = append(args, params.ProjectID)
		argIdx++
	} else {
		finalQuery = fmt.Sprintf("SELECT * FROM (%s) AS results", finalQuery)
	}

	finalQuery += fmt.Sprintf(" ORDER BY rank DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var res models.SearchResult
		err := rows.Scan(
			&res.ID,
			&res.Type,
			&res.Title,
			&res.Description,
			&res.ProjectID,
			&res.IssueID,
			&res.IssueKey,
			&res.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		results = append(results, res)
	}

	return results, nil
}
