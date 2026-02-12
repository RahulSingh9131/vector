package service

import (
	"context"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
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

// UpdateProfile updates a user's profile information and syncs to Clerk
func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	s.server.Logger.Debug().
		Str("user_id", id.String()).
		Msg("updating user profile")

	// Get the existing user to find their Clerk ID
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingUser == nil {
		return nil, errs.NewNotFoundError("User not found", false, nil)
	}

	// Sync to Clerk if the user has a Clerk ID
	if existingUser.ClerkUserID != "" {
		updateParams := &clerkUser.UpdateParams{}
		if params.FirstName != nil {
			updateParams.FirstName = params.FirstName
		}
		if params.LastName != nil {
			updateParams.LastName = params.LastName
		}

		_, clerkErr := clerkUser.Update(ctx, existingUser.ClerkUserID, updateParams)
		if clerkErr != nil {
			s.server.Logger.Warn().Err(clerkErr).
				Str("user_id", id.String()).
				Str("clerk_user_id", existingUser.ClerkUserID).
				Msg("failed to sync profile update to Clerk (continuing with local update)")
		}
	}

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
	s.server.Logger.Debug().Str("user_id", id.String()).Msg("UserService.RecordLogin starting")

	err := s.userRepo.UpdateLastLogin(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", id.String()).
			Str("ctx_err", fmt.Sprintf("%v", ctx.Err())).
			Msg("failed to record login in database")
		return fmt.Errorf("failed to record login: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user locally and deletes them from Clerk
func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	s.server.Logger.Info().
		Str("user_id", id.String()).
		Msg("deactivating user")

	// Get the existing user to find their Clerk ID
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if existingUser == nil {
		return errs.NewNotFoundError("User not found", false, nil)
	}

	// Delete from Clerk first
	if existingUser.ClerkUserID != "" {
		_, clerkErr := clerkUser.Delete(ctx, existingUser.ClerkUserID)
		if clerkErr != nil {
			s.server.Logger.Error().Err(clerkErr).
				Str("user_id", id.String()).
				Str("clerk_user_id", existingUser.ClerkUserID).
				Msg("failed to delete user from Clerk")
			return &errs.HTTPError{
				Code:    "BAD_GATEWAY",
				Message: "Failed to delete user from Clerk: " + clerkErr.Error(),
				Status:  502,
			}
		}
	}

	// Deactivate locally
	err = s.userRepo.Delete(ctx, id)
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

// ListUsers retrieves all users
func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		s.server.Logger.Error().Err(err).Msg("failed to list users")
		return nil, sqlerr.HandleError(err)
	}
	return users, nil
}

// CreateUser creates a new user in Clerk and then in the local DB
func (s *UserService) CreateUser(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	s.server.Logger.Info().Str("email", params.Email).Msg("creating user")

	// Step 1: Create in Clerk first
	skipPwReq := true
	clerkParams := &clerkUser.CreateParams{
		EmailAddresses:         &[]string{params.Email},
		FirstName:              params.FirstName,
		LastName:               params.LastName,
		SkipPasswordRequirement: &skipPwReq,
	}

	clerkUsr, err := clerkUser.Create(ctx, clerkParams)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("email", params.Email).
			Msg("failed to create user in Clerk")
		return nil, errs.NewBadRequestError("Failed to create user in Clerk: "+err.Error(), false, nil, nil, nil)
	}

	s.server.Logger.Info().
		Str("clerk_user_id", clerkUsr.ID).
		Str("email", params.Email).
		Msg("user created in Clerk")

	// Step 2: Save to local DB with the Clerk user ID
	params.ClerkUserID = clerkUsr.ID
	user, err := s.userRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("email", params.Email).
			Str("clerk_user_id", clerkUsr.ID).
			Msg("failed to create user in database (Clerk user was created)")
		return nil, sqlerr.HandleError(err)
	}

	return user, nil
}
