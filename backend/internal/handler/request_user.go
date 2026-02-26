package handler

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
