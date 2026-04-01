package models

import "github.com/google/uuid"

// SearchResultType represents the type of search result
type SearchResultType string

const (
	SearchResultTypeIssue   SearchResultType = "issue"
	SearchResultTypeComment SearchResultType = "comment"
)

// SearchResult represents a single item in a search result
type SearchResult struct {
	ID          uuid.UUID        `json:"id"`
	Type        SearchResultType `json:"type"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	// Additional metadata for UI
	ProjectID   *uuid.UUID       `json:"projectId,omitempty"`
	IssueID     *uuid.UUID       `json:"issueId,omitempty"` // For comments
	IssueKey    string           `json:"issueKey,omitempty"` // For issues
	Rank        float32          `json:"rank"`
}

// SearchParams represents the parameters for a search request
type SearchParams struct {
	Query     string           `json:"query"`
	Type      *SearchResultType `json:"type"` // Optional filter by type
	ProjectID *uuid.UUID       `json:"projectId"` // Optional filter by project
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
}
