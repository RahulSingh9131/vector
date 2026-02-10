package models

import (
	"time"

	"github.com/google/uuid"
)

// Organization represents an organization in the system
type Organization struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ClerkOrgID       string    `json:"clerkOrgId" db:"clerk_org_id"`
	Name             string    `json:"name" db:"name"`
	Slug             string    `json:"slug" db:"slug"`
	LogoURL          *string   `json:"logoUrl,omitempty" db:"logo_url"`
	SubscriptionTier string    `json:"subscriptionTier" db:"subscription_tier"`
	MaxMembers       int       `json:"maxMembers" db:"max_members"`
	MaxProjects      int       `json:"maxProjects" db:"max_projects"`
	IsActive         bool      `json:"isActive" db:"is_active"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateOrganizationParams represents the parameters for creating a new organization
type CreateOrganizationParams struct {
	ClerkOrgID       string
	Name             string
	Slug             string
	LogoURL          *string
	SubscriptionTier string
	MaxMembers       int
	MaxProjects      int
}

// UpdateOrganizationParams represents the parameters for updating an organization
type UpdateOrganizationParams struct {
	Name             *string
	Slug             *string
	LogoURL          *string
	SubscriptionTier *string
	MaxMembers       *int
	MaxProjects      *int
	IsActive         *bool
}
