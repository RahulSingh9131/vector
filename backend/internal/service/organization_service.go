package service

import (
	"context"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	"github.com/google/uuid"
)

// Allowed organization member roles
var allowedRoles = map[string]bool{
	"admin":  true,
	"member": true,
}

// OrganizationService handles business logic for organizations
type OrganizationService struct {
	server     *server.Server
	orgRepo    *repository.OrganizationRepository
	memberRepo *repository.OrganizationMemberRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(s *server.Server, repos *repository.Repositories) *OrganizationService {
	return &OrganizationService{
		server:     s,
		orgRepo:    repos.Organization,
		memberRepo: repos.OrganizationMember,
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

// UpdateSettings updates organization settings
func (s *OrganizationService) UpdateSettings(ctx context.Context, id uuid.UUID, params models.UpdateOrganizationParams) (*models.Organization, error) {
	s.server.Logger.Debug().
		Str("org_id", id.String()).
		Msg("updating organization settings")

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

// AddMember adds a user to an organization
func (s *OrganizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	// Validate role
	if !allowedRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: admin, member",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: admin, member"}},
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

	member, err := s.memberRepo.AddMember(ctx, models.CreateOrganizationMemberParams{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
	})

	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to add member to organization")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("member added to organization successfully")

	return member, nil
}

// RemoveMember removes a user from an organization
func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Msg("removing member from organization")

	err := s.memberRepo.RemoveMember(ctx, orgID, userID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to remove member from organization")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Msg("member removed from organization successfully")

	return nil
}

// UpdateMemberRole updates a member's role in an organization
func (s *OrganizationService) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	// Validate role
	if !allowedRoles[role] {
		return nil, errs.NewBadRequestError(
			"Invalid role. Must be one of: admin, member",
			true, nil,
			[]errs.FieldError{{Field: "role", Error: "must be one of: admin, member"}},
			nil,
		)
	}

	s.server.Logger.Debug().
		Str("org_id", orgID.String()).
		Str("user_id", userID.String()).
		Str("role", role).
		Msg("updating member role")

	member, err := s.memberRepo.UpdateRole(ctx, orgID, userID, role)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("org_id", orgID.String()).
			Str("user_id", userID.String()).
			Msg("failed to update member role")
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

// DeactivateOrganization soft deletes an organization
func (s *OrganizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	s.server.Logger.Info().
		Str("org_id", id.String()).
		Msg("deactivating organization")

	err := s.orgRepo.Delete(ctx, id)
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
