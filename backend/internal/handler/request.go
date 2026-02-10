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
