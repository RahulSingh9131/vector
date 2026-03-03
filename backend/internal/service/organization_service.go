package service

import (
	"context"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	clerkOrganization "github.com/clerk/clerk-sdk-go/v2/organization"
	clerkMembership "github.com/clerk/clerk-sdk-go/v2/organizationmembership"
	"github.com/google/uuid"
)

// allowedClerkRoles are the valid organization roles in Clerk's format.
var allowedClerkRoles = map[string]bool{
	"org:owner":  true,
	"org:admin":  true,
	"org:member": true,
	"org:guest":  true,
}

// OrganizationService handles business logic for organizations
type OrganizationService struct {
	server     *server.Server
	orgRepo    *repository.OrganizationRepository
	memberRepo *repository.OrganizationMemberRepository
	userRepo   *repository.UserRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(s *server.Server, repos *repository.Repositories) *OrganizationService {
	return &OrganizationService{
		server:     s,
		orgRepo:    repos.Organization,
		memberRepo: repos.OrganizationMember,
		userRepo:   repos.User,
	}
}

// GetOrCreateFromClerk syncs an organization from Clerk
func (s *OrganizationService) GetOrCreateFromClerk(ctx context.Context, clerkOrgID, name, slug string) (*models.Organization, error) {
	s.server.Logger.Debug().
		Str("clerk_org_id", clerkOrgID).
		Str("name", name).
		Msg("getting or creating organization from clerk")

	// Set default subscription tier for new organizations
	subscriptionTier := "free"
	maxMembers := 10
	maxProjects := 5

	org, err := s.orgRepo.GetOrCreate(ctx, models.CreateOrganizationParams{
		ClerkOrgID:       clerkOrgID,
		Name:             name,
		Slug:             slug,
		LogoURL:          nil,
		SubscriptionTier: subscriptionTier,
		MaxMembers:       maxMembers,
		MaxProjects:      maxProjects,
	})

	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("clerk_org_id", clerkOrgID).
			Msg("failed to get or create organization")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("clerk_org_id", clerkOrgID).
		Str("org_id", org.ID.String()).
		Msg("organization synced from clerk")

	return org, nil
}

// GetByID retrieves an organization by ID
func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", id.String()).
			Msg("failed to get organization by ID")
		return nil, sqlerr.HandleError(err)
	}

	if org == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	return org, nil
}

// GetByClerkID retrieves an organization by Clerk ID
func (s *OrganizationService) GetByClerkID(ctx context.Context, clerkOrgID string) (*models.Organization, error) {
	org, err := s.orgRepo.GetByClerkID(ctx, clerkOrgID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("clerk_org_id", clerkOrgID).
			Msg("failed to get organization by clerk ID")
		return nil, sqlerr.HandleError(err)
	}

	if org == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	return org, nil
}

// UpdateSettings updates organization settings and syncs to Clerk
func (s *OrganizationService) UpdateSettings(ctx context.Context, id uuid.UUID, params models.UpdateOrganizationParams) (*models.Organization, error) {
	s.server.Logger.Debug().
		Str("org_id", id.String()).
		Msg("updating organization settings")

	// Get the existing org to find the Clerk ID
	existingOrg, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingOrg == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Sync name/slug changes to Clerk
	if existingOrg.ClerkOrgID != "" {
		updateParams := &clerkOrganization.UpdateParams{}
		if params.Name != nil {
			updateParams.Name = params.Name
		}
		if params.Slug != nil {
			updateParams.Slug = params.Slug
		}

		_, clerkErr := clerkOrganization.Update(ctx, existingOrg.ClerkOrgID, updateParams)
		if clerkErr != nil {
			s.server.Logger.Warn().Err(clerkErr).
				Str("org_id", id.String()).
				Str("clerk_org_id", existingOrg.ClerkOrgID).
				Msg("failed to sync organization update to Clerk (continuing with local update)")
		}
	}

	org, err := s.orgRepo.Update(ctx, id, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", id.String()).
			Msg("failed to update organization settings")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", id.String()).
		Msg("organization settings updated successfully")

	return org, nil
}

// CheckMemberLimit checks if adding a new member would exceed the organization's limit
func (s *OrganizationService) CheckMemberLimit(ctx context.Context, orgID uuid.UUID) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return sqlerr.HandleError(err)
	}

	if org == nil {
		return errs.NewNotFoundError("Organization not found", false, nil)
	}

	currentCount, err := s.orgRepo.GetMemberCount(ctx, orgID)
	if err != nil {
		return sqlerr.HandleError(err)
	}

	if currentCount >= org.MaxMembers {
		return errs.NewBadRequestError(
			"Organization has reached the maximum member limit",
			true, nil, nil, nil,
		)
	}

	return nil
}

// AddMember adds a user to an organization — syncs to Clerk first, then mirrors to local DB.
func (s *OrganizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	// Validate role against allowed Clerk roles
	if !allowedClerkRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: org:owner, org:admin, org:member, org:guest",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: org:owner, org:admin, org:member, org:guest"}},
			nil,
		)
	}

	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("adding member to organization")

	// Check member limit
	if err := s.CheckMemberLimit(ctx, orgID); err != nil {
		return nil, err
	}

	// Look up the org to get its Clerk org ID
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if org == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Look up the user to get their Clerk user ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if user == nil {
		return nil, errs.NewNotFoundError("User not found", false, nil)
	}

	// Step 1: Add to Clerk (source of truth for org membership)
	_, clerkErr := clerkMembership.Create(ctx, &clerkMembership.CreateParams{
		OrganizationID: org.ClerkOrgID,
		UserID:         &user.ClerkUserID,
		Role:           &role,
	})
	if clerkErr != nil {
		s.server.Logger.Error().Err(clerkErr).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Str("clerk_org_id", org.ClerkOrgID).
			Str("clerk_user_id", user.ClerkUserID).
			Msg("failed to add member to Clerk organization")
		return nil, errs.NewBadRequestError("Failed to add member to Clerk: "+clerkErr.Error(), false, nil, nil, nil)
	}

	s.server.Logger.Info().
		Str("clerk_org_id", org.ClerkOrgID).
		Str("clerk_user_id", user.ClerkUserID).
		Msg("member added to Clerk organization")

	// Step 2: Mirror to local DB
	member, err := s.memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to add member to local DB (Clerk membership was created)")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("member added to organization successfully")

	return member, nil
}

// RemoveMember removes a user from an organization — syncs to Clerk first, then mirrors to local DB.
func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Msg("removing member from organization")

	// Look up the org to get its Clerk org ID
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if org == nil {
		return errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Look up the user to get their Clerk user ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if user == nil {
		return errs.NewNotFoundError("User not found", false, nil)
	}

	// Step 1: Remove from Clerk
	_, clerkErr := clerkMembership.Delete(ctx, &clerkMembership.DeleteParams{
		OrganizationID: org.ClerkOrgID,
		UserID:         user.ClerkUserID,
	})
	if clerkErr != nil {
		s.server.Logger.Error().Err(clerkErr).
			Str("org_id", orgID.String()).
			Str("clerk_org_id", org.ClerkOrgID).
			Str("clerk_user_id", user.ClerkUserID).
			Msg("failed to remove member from Clerk organization")
		return errs.NewBadRequestError("Failed to remove member from Clerk: "+clerkErr.Error(), false, nil, nil, nil)
	}

	s.server.Logger.Info().
		Str("clerk_org_id", org.ClerkOrgID).
		Str("clerk_user_id", user.ClerkUserID).
		Msg("member removed from Clerk organization")

	// Step 2: Remove from local DB
	if err := s.memberRepo.RemoveMember(ctx, orgID, userID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to remove member from local DB (Clerk membership was removed)")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Msg("member removed from organization successfully")

	return nil
}

// UpdateMemberRole updates a member's role in an organization — syncs to Clerk first, then mirrors to local DB.
func (s *OrganizationService) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	// Validate role against allowed Clerk roles
	if !allowedClerkRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: org:owner, org:admin, org:member, org:guest",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: org:owner, org:admin, org:member, org:guest"}},
			nil,
		)
	}

	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("updating member role")

	// Look up the org to get its Clerk org ID
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if org == nil {
		return nil, errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Look up the user to get their Clerk user ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if user == nil {
		return nil, errs.NewNotFoundError("User not found", false, nil)
	}

	// Step 1: Update role in Clerk
	_, clerkErr := clerkMembership.Update(ctx, &clerkMembership.UpdateParams{
		OrganizationID: org.ClerkOrgID,
		UserID:         user.ClerkUserID,
		Role:           &role,
	})
	if clerkErr != nil {
		s.server.Logger.Error().Err(clerkErr).
			Str("org_id", orgID.String()).
			Str("clerk_org_id", org.ClerkOrgID).
			Str("clerk_user_id", user.ClerkUserID).
			Str("role", role).
			Msg("failed to update member role in Clerk")
		return nil, errs.NewBadRequestError("Failed to update role in Clerk: "+clerkErr.Error(), false, nil, nil, nil)
	}

	s.server.Logger.Info().
		Str("clerk_org_id", org.ClerkOrgID).
		Str("clerk_user_id", user.ClerkUserID).
		Str("role", role).
		Msg("member role updated in Clerk")

	// Step 2: Mirror to local DB
	member, err := s.memberRepo.UpdateRole(ctx, orgID, userID, role)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to update member role in local DB (Clerk was updated)")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("member role updated successfully")

	return member, nil
}

// GetMembers retrieves all members of an organization
func (s *OrganizationService) GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	members, err := s.memberRepo.GetMembersByOrganization(ctx, orgID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Msg("failed to get organization members")
		return nil, sqlerr.HandleError(err)
	}

	return members, nil
}

// GetUserOrganizations retrieves all organizations a user belongs to
func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	orgs, err := s.memberRepo.GetOrganizationsByUser(ctx, userID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("user_id", userID.String()).
			Msg("failed to get user organizations")
		return nil, sqlerr.HandleError(err)
	}

	return orgs, nil
}

// DeactivateOrganization deactivates an organization locally and deletes from Clerk
func (s *OrganizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	s.server.Logger.Info().
		Str("org_id", id.String()).
		Msg("deactivating organization")

	// Get the existing org to find the Clerk ID
	existingOrg, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if existingOrg == nil {
		return errs.NewNotFoundError("Organization not found", false, nil)
	}

	// Delete from Clerk first
	if existingOrg.ClerkOrgID != "" {
		_, clerkErr := clerkOrganization.Delete(ctx, existingOrg.ClerkOrgID)
		if clerkErr != nil {
			s.server.Logger.Error().Err(clerkErr).
				Str("org_id", id.String()).
				Str("clerk_org_id", existingOrg.ClerkOrgID).
				Msg("failed to delete organization from Clerk")
			return &errs.HTTPError{
				Code:    "BAD_GATEWAY",
				Message: "Failed to delete organization from Clerk: " + clerkErr.Error(),
				Status:  502,
			}
		}
	}

	// Deactivate locally
	err = s.orgRepo.Delete(ctx, id)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", id.String()).
			Msg("failed to deactivate organization")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", id.String()).
		Msg("organization deactivated successfully")

	return nil
}

// CreateOrganization creates a new organization in Clerk and then in the local DB
func (s *OrganizationService) CreateOrganization(ctx context.Context, params models.CreateOrganizationParams) (*models.Organization, error) {
	s.server.Logger.Info().Str("name", params.Name).Msg("creating organization")

	// Set defaults if not provided
	if params.SubscriptionTier == "" {
		params.SubscriptionTier = "free"
	}
	if params.MaxMembers == 0 {
		params.MaxMembers = 10
	}
	if params.MaxProjects == 0 {
		params.MaxProjects = 5
	}

	// Step 1: Create in Clerk first so it appears in the dashboard
	clerkOrg, err := clerkOrganization.Create(ctx, &clerkOrganization.CreateParams{
		Name: &params.Name,
		Slug: &params.Slug,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("name", params.Name).
			Msg("failed to create organization in Clerk")
		return nil, errs.NewBadRequestError("Failed to create organization in Clerk: "+err.Error(), false, nil, nil, nil)
	}

	s.server.Logger.Info().
		Str("clerk_org_id", clerkOrg.ID).
		Str("name", params.Name).
		Msg("organization created in Clerk")

	// Step 2: Save to local DB with the Clerk org ID
	params.ClerkOrgID = clerkOrg.ID
	org, err := s.orgRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("name", params.Name).
			Str("clerk_org_id", clerkOrg.ID).
			Msg("failed to create organization in database (Clerk org was created)")
		return nil, sqlerr.HandleError(err)
	}

	return org, nil
}
