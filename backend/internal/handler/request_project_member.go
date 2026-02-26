package handler

// ListProjectMembersRequest represents path params for listing project members
type ListProjectMembersRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
}

func (r *ListProjectMembersRequest) Validate() error {
	return validate.Struct(r)
}

// AddProjectMemberRequest represents the request body for adding a project member
type AddProjectMemberRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `json:"user_id" validate:"required,uuid"`
	Role      string `json:"role" validate:"required,oneof=admin member viewer"`
}

func (r *AddProjectMemberRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateProjectMemberRoleRequest represents the request body for updating a project member's role
type UpdateProjectMemberRoleRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `param:"userId" validate:"required,uuid"`
	Role      string `json:"role" validate:"required,oneof=admin member viewer"`
}

func (r *UpdateProjectMemberRoleRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveProjectMemberRequest represents path params for removing a project member
type RemoveProjectMemberRequest struct {
	OrgID     string `param:"orgId" validate:"required,uuid"`
	ProjectID string `param:"projectId" validate:"required,uuid"`
	UserID    string `param:"userId" validate:"required,uuid"`
}

func (r *RemoveProjectMemberRequest) Validate() error {
	return validate.Struct(r)
}
