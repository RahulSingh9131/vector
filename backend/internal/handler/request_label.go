package handler

// CreateLabelRequest represents the request body for creating a label
type CreateLabelRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Color       string  `json:"color" validate:"required,min=7,max=7"`
	Description *string `json:"description" validate:"omitempty"`
}

func (r *CreateLabelRequest) Validate() error {
	return validate.Struct(r)
}

// ListLabelsRequest represents path params for listing labels
type ListLabelsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *ListLabelsRequest) Validate() error {
	return validate.Struct(r)
}

// GetLabelRequest represents path params for getting a label
type GetLabelRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *GetLabelRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateLabelRequest represents the request body for updating a label
type UpdateLabelRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	LabelID     string  `param:"labelId" validate:"required,uuid"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Color       *string `json:"color" validate:"omitempty,min=7,max=7"`
	Description *string `json:"description" validate:"omitempty"`
}

func (r *UpdateLabelRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteLabelRequest represents path params for deleting a label
type DeleteLabelRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *DeleteLabelRequest) Validate() error {
	return validate.Struct(r)
}

// GetIssueLabelsRequest represents path params for getting labels on an issue
type GetIssueLabelsRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
}

func (r *GetIssueLabelsRequest) Validate() error {
	return validate.Struct(r)
}

// AddLabelToIssueRequest represents the request body for adding a label to an issue
type AddLabelToIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	LabelID   string `json:"label_id" validate:"required,uuid"`
}

func (r *AddLabelToIssueRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveLabelFromIssueRequest represents path params for removing a label from an issue
type RemoveLabelFromIssueRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	IssueID   string `param:"issueId" validate:"required,uuid"`
	LabelID   string `param:"labelId" validate:"required,uuid"`
}

func (r *RemoveLabelFromIssueRequest) Validate() error {
	return validate.Struct(r)
}
