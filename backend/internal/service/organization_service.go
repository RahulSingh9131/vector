package service

import (
	"context"
	"fmt"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
)

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
		return nil, fmt.Errorf("failed to get or create organization: %w", err)
	}

	return org, nil
}

// GetByID retrieves an organization by ID
func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	return org, nil
}

// GetByClerkID retrieves an organization by Clerk ID
func (s *OrganizationService) GetByClerkID(ctx context.Context, clerkOrgID string) (*models.Organization, error) {
	org, err := s.orgRepo.GetByClerkID(ctx, clerkOrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	return org, nil
}

// UpdateSettings updates organization settings
func (s *OrganizationService) UpdateSettings(ctx context.Context, id uuid.UUID, params models.UpdateOrganizationParams) (*models.Organization, error) {
	org, err := s.orgRepo.Update(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return org, nil
}

// CheckMemberLimit checks if adding a new member would exceed the organization's limit
func (s *OrganizationService) CheckMemberLimit(ctx context.Context, orgID uuid.UUID) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}

	if org == nil {
		return fmt.Errorf("organization not found")
	}

	currentCount, err := s.orgRepo.GetMemberCount(ctx, orgID)
	if err != nil {
		return fmt.Errorf("failed to get member count: %w", err)
	}

	if currentCount >= org.MaxMembers {
		return fmt.Errorf("organization has reached maximum member limit of %d", org.MaxMembers)
	}

	return nil
}

// AddMember adds a user to an organization
func (s *OrganizationService) AddMember(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
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
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return member, nil
}

// RemoveMember removes a user from an organization
func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	err := s.memberRepo.RemoveMember(ctx, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

// UpdateMemberRole updates a member's role in an organization
func (s *OrganizationService) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	member, err := s.memberRepo.UpdateRole(ctx, orgID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to update member role: %w", err)
	}

	return member, nil
}

// GetMembers retrieves all members of an organization
func (s *OrganizationService) GetMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	members, err := s.memberRepo.GetMembersByOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	return members, nil
}

// GetUserOrganizations retrieves all organizations a user belongs to
func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	orgs, err := s.memberRepo.GetOrganizationsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	return orgs, nil
}

// DeactivateOrganization soft deletes an organization
func (s *OrganizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	err := s.orgRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate organization: %w", err)
	}

	return nil
}
