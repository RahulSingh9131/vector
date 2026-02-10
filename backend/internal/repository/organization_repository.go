package repository

import (
	"context"
	"errors"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// OrganizationRepository handles database operations for organizations
type OrganizationRepository struct {
	server *server.Server
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(s *server.Server) *OrganizationRepository {
	return &OrganizationRepository{
		server: s,
	}
}

// GetByID retrieves an organization by their internal UUID
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	query := `
		SELECT id, clerk_org_id, name, slug, logo_url, subscription_tier, 
		       max_members, max_projects, is_active, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org models.Organization
	err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(
		&org.ID,
		&org.ClerkOrgID,
		&org.Name,
		&org.Slug,
		&org.LogoURL,
		&org.SubscriptionTier,
		&org.MaxMembers,
		&org.MaxProjects,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

// GetByClerkID retrieves an organization by their Clerk organization ID
func (r *OrganizationRepository) GetByClerkID(ctx context.Context, clerkOrgID string) (*models.Organization, error) {
	query := `
		SELECT id, clerk_org_id, name, slug, logo_url, subscription_tier, 
		       max_members, max_projects, is_active, created_at, updated_at
		FROM organizations
		WHERE clerk_org_id = $1
	`

	var org models.Organization
	err := r.server.DB.Pool.QueryRow(ctx, query, clerkOrgID).Scan(
		&org.ID,
		&org.ClerkOrgID,
		&org.Name,
		&org.Slug,
		&org.LogoURL,
		&org.SubscriptionTier,
		&org.MaxMembers,
		&org.MaxProjects,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

// GetBySlug retrieves an organization by their slug
func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	query := `
		SELECT id, clerk_org_id, name, slug, logo_url, subscription_tier, 
		       max_members, max_projects, is_active, created_at, updated_at
		FROM organizations
		WHERE slug = $1
	`

	var org models.Organization
	err := r.server.DB.Pool.QueryRow(ctx, query, slug).Scan(
		&org.ID,
		&org.ClerkOrgID,
		&org.Name,
		&org.Slug,
		&org.LogoURL,
		&org.SubscriptionTier,
		&org.MaxMembers,
		&org.MaxProjects,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, params models.CreateOrganizationParams) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (clerk_org_id, name, slug, logo_url, subscription_tier, max_members, max_projects)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, clerk_org_id, name, slug, logo_url, subscription_tier, 
		          max_members, max_projects, is_active, created_at, updated_at
	`

	var org models.Organization
	err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		params.ClerkOrgID,
		params.Name,
		params.Slug,
		params.LogoURL,
		params.SubscriptionTier,
		params.MaxMembers,
		params.MaxProjects,
	).Scan(
		&org.ID,
		&org.ClerkOrgID,
		&org.Name,
		&org.Slug,
		&org.LogoURL,
		&org.SubscriptionTier,
		&org.MaxMembers,
		&org.MaxProjects,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &org, nil
}

// Update updates an organization's information
func (r *OrganizationRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateOrganizationParams) (*models.Organization, error) {
	query := `
		UPDATE organizations
		SET name = COALESCE($2, name),
		    slug = COALESCE($3, slug),
		    logo_url = COALESCE($4, logo_url),
		    subscription_tier = COALESCE($5, subscription_tier),
		    max_members = COALESCE($6, max_members),
		    max_projects = COALESCE($7, max_projects),
		    is_active = COALESCE($8, is_active)
		WHERE id = $1
		RETURNING id, clerk_org_id, name, slug, logo_url, subscription_tier, 
		          max_members, max_projects, is_active, created_at, updated_at
	`

	var org models.Organization
	err := r.server.DB.Pool.QueryRow(
		ctx,
		query,
		id,
		params.Name,
		params.Slug,
		params.LogoURL,
		params.SubscriptionTier,
		params.MaxMembers,
		params.MaxProjects,
		params.IsActive,
	).Scan(
		&org.ID,
		&org.ClerkOrgID,
		&org.Name,
		&org.Slug,
		&org.LogoURL,
		&org.SubscriptionTier,
		&org.MaxMembers,
		&org.MaxProjects,
		&org.IsActive,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &org, nil
}

// Delete soft deletes an organization by setting is_active to false
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE organizations
		SET is_active = false
		WHERE id = $1
	`

	_, err := r.server.DB.Pool.Exec(ctx, query, id)
	return err
}

// GetOrCreate retrieves an organization by Clerk ID or creates it if it doesn't exist (atomic)
func (r *OrganizationRepository) GetOrCreate(ctx context.Context, params models.CreateOrganizationParams) (*models.Organization, error) {
	query := `
		INSERT INTO organizations (clerk_org_id, name, slug, logo_url, subscription_tier, max_members, max_projects)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (clerk_org_id) DO NOTHING
	`

	_, err := r.server.DB.Pool.Exec(
		ctx,
		query,
		params.ClerkOrgID,
		params.Name,
		params.Slug,
		params.LogoURL,
		params.SubscriptionTier,
		params.MaxMembers,
		params.MaxProjects,
	)
	if err != nil {
		return nil, err
	}

	// Always fetch the organization (whether just inserted or already existed)
	return r.GetByClerkID(ctx, params.ClerkOrgID)
}

// GetMemberCount returns the number of members in an organization
func (r *OrganizationRepository) GetMemberCount(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM organization_members
		WHERE organization_id = $1
	`

	var count int
	err := r.server.DB.Pool.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
