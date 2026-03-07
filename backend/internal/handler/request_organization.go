package handler

// GetOrganizationRequest represents path params for getting an organization by ID
type GetOrganizationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetOrganizationRequest) Validate() error {
	return validate.Struct(r)
}

// ListMembersRequest represents path params for listing members of an organization
type ListMembersRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *ListMembersRequest) Validate() error {
	return validate.Struct(r)
}

// CreateOrganizationRequest represents the request body for creating an organization
type CreateOrganizationRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Slug    string  `json:"slug" validate:"required,min=1,max=100"`
	LogoURL *string `json:"logoUrl" validate:"omitempty,url"`
}

func (r *CreateOrganizationRequest) Validate() error {
	return validate.Struct(r)
}

// AddMemberRequest represents the request body for adding a member to an organization
type AddMemberRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=org:owner org:admin org:member org:guest"`
}

func (r *AddMemberRequest) Validate() error {
	return validate.Struct(r)
}

// UpdateMemberRoleRequest represents the request body for updating a member's role
type UpdateMemberRoleRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `param:"userId" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=admin member guest"`
}

func (r *UpdateMemberRoleRequest) Validate() error {
	return validate.Struct(r)
}

// RemoveMemberRequest represents path params for removing a member
type RemoveMemberRequest struct {
	ID     string `param:"id" validate:"required,uuid"`
	UserID string `param:"userId" validate:"required,uuid"`
}

func (r *RemoveMemberRequest) Validate() error {
	return validate.Struct(r)
}
