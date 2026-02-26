package handler

// CreateIssueRequest represents the request body for creating an issue
type CreateIssueRequest struct {
	OrgID         string  `param:"orgId" validate:"required,uuid"`
	ProjectID     string  `param:"projectId" validate:"required,uuid"`
	Title         string  `json:"title" validate:"required,min=1,max=500"`
	Description   *string `json:"description" validate:"omitempty"`
	Priority      string  `json:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type          string  `json:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID    *string `json:"assignee_id" validate:"omitempty,uuid"`
	ParentIssueID *string `json:"parent_issue_id" validate:"omitempty,uuid"`
	DueDate       *string `json:"due_date" validate:"omitempty"`
}

func (r *CreateIssueRequest) Validate() error {
	return validate.Struct(r)
}

// ListIssuesRequest represents query params for listing issues
type ListIssuesRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	Status     *string `query:"status" validate:"omitempty,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
	Priority   *string `query:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type       *string `query:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID *string `query:"assignee_id" validate:"omitempty,uuid"`
	Page       int     `query:"page" validate:"omitempty,min=1"`
	Limit      int     `query:"limit" validate:"omitempty,min=1,max=100"`
}

func (r *ListIssuesRequest) Validate() error {
	return validate.Struct(r)
}

// GetIssueRequest represents path params for getting an issue
type GetIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *GetIssueRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateIssueRequest represents the request body for updating an issue
type UpdateIssueRequest struct {
	OrgID         string  `param:"orgId" validate:"required,uuid"`
	ProjectID     string  `param:"projectId" validate:"required,uuid"`
	IssueID       string  `param:"issueId" validate:"required,uuid"`
	Title         *string `json:"title" validate:"omitempty,min=1,max=500"`
	Description   *string `json:"description" validate:"omitempty"`
	Status        *string `json:"status" validate:"omitempty,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
	Priority      *string `json:"priority" validate:"omitempty,oneof=urgent high medium low none"`
	Type          *string `json:"type" validate:"omitempty,oneof=task bug story epic"`
	AssigneeID    *string `json:"assignee_id" validate:"omitempty,uuid"`
	SortOrder     *int    `json:"sort_order" validate:"omitempty,min=0"`
	ParentIssueID *string `json:"parent_issue_id" validate:"omitempty,uuid"`
	DueDate       *string `json:"due_date" validate:"omitempty"`
}

func (r *UpdateIssueRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteIssueRequest represents path params for deleting an issue
type DeleteIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *DeleteIssueRequest) Validate() error {
	return validate.Struct(r)
}

// AssignIssueRequest represents the request body for assigning an issue
type AssignIssueRequest struct {
	OrgID      string  `param:"orgId" validate:"required,uuid"`
	ProjectID  string  `param:"projectId" validate:"required,uuid"`
	IssueID    string  `param:"issueId" validate:"required,uuid"`
	AssigneeID *string `json:"assignee_id" validate:"omitempty,uuid"`
}

func (r *AssignIssueRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateIssueStatusRequest represents the request body for updating an issue's status
type UpdateIssueStatusRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	Status    string `json:"status" validate:"required,oneof=backlog todo in_progress in_review in_dev in_prod cancelled"`
}

func (r *UpdateIssueStatusRequest) Validate() error {
	return validate.Struct(r)
}
