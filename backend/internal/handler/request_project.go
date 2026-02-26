package handler

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Slug        string  `json:"slug" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
	Identifier  string  `json:"identifier" validate:"required,min=2,max=10,uppercase"`
}

func (r *CreateProjectRequest) Validate() error {
	return validate.Struct(r)
}

// ListProjectsRequest represents path params for listing projects
type ListProjectsRequest struct {
	OrgID string `param:"orgId" validate:"required,uuid"`
}

func (r *ListProjectsRequest) Validate() error {
	return validate.Struct(r)
}

// GetProjectRequest represents path params for getting a project
type GetProjectRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *GetProjectRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateProjectRequest represents the request body for updating a project
type UpdateProjectRequest struct {
	OrgID       string  `param:"orgId" validate:"required,uuid"`
	ProjectID   string  `param:"projectId" validate:"required,uuid"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Slug        *string `json:"slug" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=active archived"`
}

func (r *UpdateProjectRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteProjectRequest represents path params for deleting a project
type DeleteProjectRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *DeleteProjectRequest) Validate() error {
	return validate.Struct(r)
}
