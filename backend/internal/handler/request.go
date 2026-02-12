package handler

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// --- User Request Structs ---

// EmptyRequest is a no-op request for endpoints that don't accept input
type EmptyRequest struct{}

func (r *EmptyRequest) Validate() error {
	return nil
}

// EmptyResponse is a no-op response for endpoints that don't return data
type EmptyResponse struct{}

// UpdateCurrentUserRequest represents the request body for updating the current user's profile
type UpdateCurrentUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName  *string `json:"lastName" validate:"omitempty,min=1,max=100"`
	AvatarURL *string `json:"avatarUrl" validate:"omitempty,url"`
}

func (r *UpdateCurrentUserRequest) Validate() error {
	return validate.Struct(r)
}

// CreateUserRequest represents the request body for manually creating a user
type CreateUserRequest struct {
	ClerkUserID string  `json:"clerk_user_id" validate:"required"`
	Email       string  `json:"email" validate:"required,email"`
	FirstName   *string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName    *string `json:"lastName" validate:"omitempty,min=1,max=100"`
	AvatarURL   *string `json:"avatarUrl" validate:"omitempty,url"`
}

func (r *CreateUserRequest) Validate() error {
	return validate.Struct(r)
}

// GetUserRequest represents path params for getting a user by ID
type GetUserRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetUserRequest) Validate() error {
	return validate.Struct(r)
}

// DeleteUserRequest represents path params for deleting a user
type DeleteUserRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteUserRequest) Validate() error {
	return validate.Struct(r)
}

// --- Organization Request Structs ---

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
	Role   string `json:"role" validate:"required,oneof=admin member guest"`
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
