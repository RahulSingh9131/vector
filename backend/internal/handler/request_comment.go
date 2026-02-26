package handler

// CreateCommentRequest represents the request body for creating a comment
type CreateCommentRequest struct {
	OrgID           string  `param:"orgId" validate:"required,uuid"`
	ProjectID       string  `param:"projectId" validate:"required,uuid"`
	IssueID         string  `param:"issueId" validate:"required,uuid"`
	Body            string  `json:"body" validate:"required,min=1"`
	ParentCommentID *string `json:"parent_comment_id" validate:"omitempty,uuid"`
}

func (r *CreateCommentRequest) Validate() error {
	return validate.Struct(r)
}

// ListCommentsRequest represents query params for listing comments
type ListCommentsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	Page      int    `query:"page" validate:"omitempty,min=1"`
	Limit     int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *ListCommentsRequest) Validate() error {
	return validate.Struct(r)
}

// GetCommentRequest represents path params for getting a comment
type GetCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
}

func (r *GetCommentRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateCommentRequest represents the request body for updating a comment
type UpdateCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
	Body      string `json:"body" validate:"required,min=1"`
}

func (r *UpdateCommentRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteCommentRequest represents path params for deleting a comment
type DeleteCommentRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	CommentID string `param:"commentId" validate:"required,uuid"`
}

func (r *DeleteCommentRequest) Validate() error {
	return validate.Struct(r)
}
