package service

import (
	"context"
	"fmt"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

// UserService handles business logic for users
type UserService struct {
	server   *server.Server
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(s *server.Server, repos *repository.Repositories) *UserService {
	return &UserService{
		server:   s,
		userRepo: repos.User,
	}
}

// GetOrCreateFromClerk syncs a user from Clerk session claims
// For full user details, use SyncUserFromClerk which fetches from Clerk API
func (s *UserService) GetOrCreateFromClerk(ctx context.Context, clerkUserID string) (*models.User, error) {
	// Try to get existing user
	user, err := s.userRepo.GetByClerkID(ctx, clerkUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// If user exists, return it
	if user != nil {
		return user, nil
	}

	// Create minimal user record - will be enriched via webhook or API call
	user, err = s.userRepo.Create(ctx, models.CreateUserParams{
		ClerkUserID: clerkUserID,
		Email:       clerkUserID + "@temp.clerk", // Temporary email, will be updated
		FirstName:   nil,
		LastName:    nil,
		AvatarURL:   nil,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// TODO: Trigger async job to fetch full user details from Clerk API
	s.server.Logger.Info().
		Str("clerk_user_id", clerkUserID).
		Str("user_id", user.ID.String()).
		Msg("created minimal user record, full sync pending")

	return user, nil
}

// GetByID retrieves a user by their ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetByClerkID retrieves a user by their Clerk ID
func (s *UserService) GetByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	user, err := s.userRepo.GetByClerkID(ctx, clerkUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// UpdateProfile updates a user's profile information
func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	user, err := s.userRepo.Update(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// RecordLogin updates the user's last login timestamp
func (s *UserService) RecordLogin(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.UpdateLastLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to record login: %w", err)
	}

	return nil
}

// DeactivateUser soft deletes a user
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}
