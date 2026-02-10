package service

import (
	"context"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
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
func (s *UserService) GetOrCreateFromClerk(ctx context.Context, clerkUserID string) (*models.User, error) {
	s.server.Logger.Debug().
		Str("clerk_user_id", clerkUserID).
		Msg("getting or creating user from clerk")

	user, err := s.userRepo.GetOrCreate(ctx, models.CreateUserParams{
		ClerkUserID: clerkUserID,
		Email:       clerkUserID + "@temp.clerk", // Temporary email, will be updated via webhook
		FirstName:   nil,
		LastName:    nil,
		AvatarURL:   nil,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("clerk_user_id", clerkUserID).
			Msg("failed to get or create user from clerk")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("clerk_user_id", clerkUserID).
		Str("user_id", user.ID.String()).
		Msg("user synced from clerk")

	return user, nil
}

// GetByID retrieves a user by their ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", id.String()).
			Msg("failed to get user by ID")
		return nil, sqlerr.HandleError(err)
	}

	if user == nil {
		return nil, errs.NewNotFoundError("User not found", false, nil)
	}

	return user, nil
}

// GetByClerkID retrieves a user by their Clerk ID
func (s *UserService) GetByClerkID(ctx context.Context, clerkUserID string) (*models.User, error) {
	user, err := s.userRepo.GetByClerkID(ctx, clerkUserID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("clerk_user_id", clerkUserID).
			Msg("failed to get user by clerk ID")
		return nil, sqlerr.HandleError(err)
	}

	if user == nil {
		return nil, errs.NewNotFoundError("User not found", false, nil)
	}

	return user, nil
}

// UpdateProfile updates a user's profile information
func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	s.server.Logger.Debug().
		Str("user_id", id.String()).
		Msg("updating user profile")

	user, err := s.userRepo.Update(ctx, id, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", id.String()).
			Msg("failed to update user profile")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("user_id", id.String()).
		Msg("user profile updated successfully")

	return user, nil
}

// RecordLogin updates the user's last login timestamp
func (s *UserService) RecordLogin(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.UpdateLastLogin(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", id.String()).
			Msg("failed to record login")
		return fmt.Errorf("failed to record login: %w", err)
	}

	return nil
}

// DeactivateUser soft deletes a user
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	s.server.Logger.Info().
		Str("user_id", id.String()).
		Msg("deactivating user")

	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", id.String()).
			Msg("failed to deactivate user")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("user_id", id.String()).
		Msg("user deactivated successfully")

	return nil
}
