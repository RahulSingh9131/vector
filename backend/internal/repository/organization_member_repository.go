package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// OrganizationMemberRepository handles database operations for organization members
type OrganizationMemberRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewOrganizationMemberRepository creates a new organization member repository
func NewOrganizationMemberRepository(s *server.Server) *OrganizationMemberRepository {
	return &OrganizationMemberRepository{
		db:     s.DB.Pool,
		server: s,
	}
}

// AddMember adds a user to an organization with a specific role
func (r *OrganizationMemberRepository) AddMember(ctx context.Context, params models.CreateOrganizationMemberParams) (*models.OrganizationMember, error) {
	query := `
		INSERT INTO organization_members (organization_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (organization_id, user_id) 
		DO UPDATE SET role = EXCLUDED.role
		RETURNING id, organization_id, user_id, role, joined_at
	`

	var member models.OrganizationMember
	err := r.db.QueryRow(
		ctx,
		query,
		params.OrganizationID,
		params.UserID,
		params.Role,
	).Scan(
		&member.ID,
		&member.OrganizationID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		return nil, err
	}

	return &member, nil
}

// RemoveMember removes a user from an organization
func (r *OrganizationMemberRepository) RemoveMember(ctx context.Context, organizationID, userID uuid.UUID) error {
	query := `
		DELETE FROM organization_members
		WHERE organization_id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, organizationID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("member not found in organization")
	}

	return nil
}

// UpdateRole updates a member's role in an organization
func (r *OrganizationMemberRepository) UpdateRole(ctx context.Context, organizationID, userID uuid.UUID, role string) (*models.OrganizationMember, error) {
	query := `
		UPDATE organization_members
		SET role = $3
		WHERE organization_id = $1 AND user_id = $2
		RETURNING id, organization_id, user_id, role, joined_at
	`

	var member models.OrganizationMember
	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
		userID,
		role,
	).Scan(
		&member.ID,
		&member.OrganizationID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("member not found in organization")
		}
		return nil, err
	}

	return &member, nil
}

// GetMember retrieves a specific member's information
func (r *OrganizationMemberRepository) GetMember(ctx context.Context, organizationID, userID uuid.UUID) (*models.OrganizationMember, error) {
	query := `
		SELECT id, organization_id, user_id, role, joined_at
		FROM organization_members
		WHERE organization_id = $1 AND user_id = $2
	`

	var member models.OrganizationMember
	err := r.db.QueryRow(ctx, query, organizationID, userID).Scan(
		&member.ID,
		&member.OrganizationID,
		&member.UserID,
		&member.Role,
		&member.JoinedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}

// GetMembersByOrganization retrieves all members of an organization with user details
func (r *OrganizationMemberRepository) GetMembersByOrganization(ctx context.Context, organizationID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	query := `
		SELECT 
			om.id, om.organization_id, om.user_id, om.role, om.joined_at,
			u.id, u.clerk_user_id, u.email, u.first_name, u.last_name, 
			u.avatar_url, u.is_active, u.last_login_at, u.created_at, u.updated_at
		FROM organization_members om
		INNER JOIN users u ON om.user_id = u.id
		WHERE om.organization_id = $1
		ORDER BY om.joined_at ASC
	`

	rows, err := r.db.Query(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.OrganizationMemberWithDetails
	for rows.Next() {
		var member models.OrganizationMemberWithDetails
		member.User = &models.User{}

		err := rows.Scan(
			&member.ID,
			&member.OrganizationID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.User.ID,
			&member.User.ClerkUserID,
			&member.User.Email,
			&member.User.FirstName,
			&member.User.LastName,
			&member.User.AvatarURL,
			&member.User.IsActive,
			&member.User.LastLoginAt,
			&member.User.CreatedAt,
			&member.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// GetOrganizationsByUser retrieves all organizations a user belongs to with organization details
func (r *OrganizationMemberRepository) GetOrganizationsByUser(ctx context.Context, userID uuid.UUID) ([]models.OrganizationMemberWithDetails, error) {
	query := `
		SELECT 
			om.id, om.organization_id, om.user_id, om.role, om.joined_at,
			o.id, o.clerk_org_id, o.name, o.slug, o.logo_url, o.subscription_tier,
			o.max_members, o.max_projects, o.is_active, o.created_at, o.updated_at
		FROM organization_members om
		INNER JOIN organizations o ON om.organization_id = o.id
		WHERE om.user_id = $1
		ORDER BY om.joined_at ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []models.OrganizationMemberWithDetails
	for rows.Next() {
		var membership models.OrganizationMemberWithDetails
		membership.Organization = &models.Organization{}

		err := rows.Scan(
			&membership.ID,
			&membership.OrganizationID,
			&membership.UserID,
			&membership.Role,
			&membership.JoinedAt,
			&membership.Organization.ID,
			&membership.Organization.ClerkOrgID,
			&membership.Organization.Name,
			&membership.Organization.Slug,
			&membership.Organization.LogoURL,
			&membership.Organization.SubscriptionTier,
			&membership.Organization.MaxMembers,
			&membership.Organization.MaxProjects,
			&membership.Organization.IsActive,
			&membership.Organization.CreatedAt,
			&membership.Organization.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

// GetMembersByRole retrieves all members of an organization with a specific role
func (r *OrganizationMemberRepository) GetMembersByRole(ctx context.Context, organizationID uuid.UUID, role string) ([]models.OrganizationMemberWithDetails, error) {
	query := `
		SELECT 
			om.id, om.organization_id, om.user_id, om.role, om.joined_at,
			u.id, u.clerk_user_id, u.email, u.first_name, u.last_name, 
			u.avatar_url, u.is_active, u.last_login_at, u.created_at, u.updated_at
		FROM organization_members om
		INNER JOIN users u ON om.user_id = u.id
		WHERE om.organization_id = $1 AND om.role = $2
		ORDER BY om.joined_at ASC
	`

	rows, err := r.db.Query(ctx, query, organizationID, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.OrganizationMemberWithDetails
	for rows.Next() {
		var member models.OrganizationMemberWithDetails
		member.User = &models.User{}

		err := rows.Scan(
			&member.ID,
			&member.OrganizationID,
			&member.UserID,
			&member.Role,
			&member.JoinedAt,
			&member.User.ID,
			&member.User.ClerkUserID,
			&member.User.Email,
			&member.User.FirstName,
			&member.User.LastName,
			&member.User.AvatarURL,
			&member.User.IsActive,
			&member.User.LastLoginAt,
			&member.User.CreatedAt,
			&member.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
