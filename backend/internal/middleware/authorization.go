package middleware

import (
	"github.com/RahulSingh9131/vector/internal/errs"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AuthorizationMiddleware provides RBAC enforcement at the router level.
// It verifies organization membership, organization roles, and project roles.
type AuthorizationMiddleware struct {
	server  *server.Server
	orgRepo *repository.OrganizationMemberRepository
	prjRepo *repository.ProjectMemberRepository
}

// NewAuthorizationMiddleware creates a new authorization middleware instance.
func NewAuthorizationMiddleware(s *server.Server, repos *repository.Repositories) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		server:  s,
		orgRepo: repos.OrganizationMember,
		prjRepo: repos.ProjectMember,
	}
}

// RequireOrgMember returns middleware that verifies the authenticated user is a member
// of the organization specified by the :id or :orgId URL parameter.
//
// If allowedRoles is provided, it additionally enforces that the user's org role
// is one of the specified roles (e.g. "admin", "member", "guest").
// If no roles are specified, any org member is allowed through.
//
// On success, sets "org_role" in the echo context for downstream use.
func (a *AuthorizationMiddleware) RequireOrgMember(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract user ID from context (set by RequireAuth)
			userID, ok := c.Get("user_id").(uuid.UUID)
			if !ok {
				a.server.Logger.Error().
					Str("middleware", "RequireOrgMember").
					Str("request_id", GetRequestID(c)).
					Msg("user_id missing from context — RequireAuth must run first")
				return errs.NewUnauthorizedError("Authentication required", false)
			}

			// Extract org ID from URL params — supports both :id and :orgId naming
			orgIDStr := c.Param("orgId")
			if orgIDStr == "" {
				orgIDStr = c.Param("id")
			}
			if orgIDStr == "" {
				a.server.Logger.Error().
					Str("middleware", "RequireOrgMember").
					Str("request_id", GetRequestID(c)).
					Msg("organization ID not found in URL parameters")
				return errs.NewBadRequestError("Organization ID is required", false, nil, nil, nil)
			}

			orgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				a.server.Logger.Warn().
					Str("middleware", "RequireOrgMember").
					Str("org_id_raw", orgIDStr).
					Str("request_id", GetRequestID(c)).
					Msg("invalid organization ID format")
				return errs.NewBadRequestError("Invalid organization ID format", false, nil, nil, nil)
			}

			// Look up membership
			member, err := a.orgRepo.GetMember(c.Request().Context(), orgID, userID)
			if err != nil {
				a.server.Logger.Error().
					Err(err).
					Str("middleware", "RequireOrgMember").
					Str("org_id", orgID.String()).
					Str("user_id", userID.String()).
					Str("request_id", GetRequestID(c)).
					Msg("failed to query organization membership")
				return errs.NewInternalServerError()
			}

			if member == nil {
				a.server.Logger.Warn().
					Str("middleware", "RequireOrgMember").
					Str("org_id", orgID.String()).
					Str("user_id", userID.String()).
					Str("request_id", GetRequestID(c)).
					Msg("access denied — user is not a member of this organization")
				return errs.NewForbiddenError("You are not a member of this organization", false)
			}

			// Check role if specific roles are required
			if len(allowedRoles) > 0 {
				if !containsRole(allowedRoles, member.Role) {
					a.server.Logger.Warn().
						Str("middleware", "RequireOrgMember").
						Str("org_id", orgID.String()).
						Str("user_id", userID.String()).
						Str("user_role", member.Role).
						Strs("required_roles", allowedRoles).
						Str("request_id", GetRequestID(c)).
						Msg("access denied — insufficient organization role")
					return errs.NewForbiddenError("You do not have the required role to perform this action", false)
				}
			}

			// Set org role in context for downstream handlers
			c.Set("org_role", member.Role)

			return next(c)
		}
	}
}

// RequireProjectRole returns middleware that verifies the authenticated user is a member
// of the project specified by the :projectId URL parameter.
//
// If allowedRoles is provided, it additionally enforces that the user's project role
// is one of the specified roles (e.g. "admin", "member", "viewer").
// If no roles are specified, any project member is allowed through.
//
// On success, sets "project_role" in the echo context for downstream use.
func (a *AuthorizationMiddleware) RequireProjectRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract user ID from context (set by RequireAuth)
			userID, ok := c.Get("user_id").(uuid.UUID)
			if !ok {
				a.server.Logger.Error().
					Str("middleware", "RequireProjectRole").
					Str("request_id", GetRequestID(c)).
					Msg("user_id missing from context — RequireAuth must run first")
				return errs.NewUnauthorizedError("Authentication required", false)
			}

			// Extract project ID from URL params
			projectIDStr := c.Param("projectId")
			if projectIDStr == "" {
				a.server.Logger.Error().
					Str("middleware", "RequireProjectRole").
					Str("request_id", GetRequestID(c)).
					Msg("project ID not found in URL parameters")
				return errs.NewBadRequestError("Project ID is required", false, nil, nil, nil)
			}

			projectID, err := uuid.Parse(projectIDStr)
			if err != nil {
				a.server.Logger.Warn().
					Str("middleware", "RequireProjectRole").
					Str("project_id_raw", projectIDStr).
					Str("request_id", GetRequestID(c)).
					Msg("invalid project ID format")
				return errs.NewBadRequestError("Invalid project ID format", false, nil, nil, nil)
			}

			// Look up membership
			member, err := a.prjRepo.GetMember(c.Request().Context(), projectID, userID)
			if err != nil {
				a.server.Logger.Error().
					Err(err).
					Str("middleware", "RequireProjectRole").
					Str("project_id", projectID.String()).
					Str("user_id", userID.String()).
					Str("request_id", GetRequestID(c)).
					Msg("failed to query project membership")
				return errs.NewInternalServerError()
			}

			if member == nil {
				a.server.Logger.Warn().
					Str("middleware", "RequireProjectRole").
					Str("project_id", projectID.String()).
					Str("user_id", userID.String()).
					Str("request_id", GetRequestID(c)).
					Msg("access denied — user is not a member of this project")
				return errs.NewForbiddenError("You are not a member of this project", false)
			}

			// Check role if specific roles are required
			if len(allowedRoles) > 0 {
				if !containsRole(allowedRoles, member.Role) {
					a.server.Logger.Warn().
						Str("middleware", "RequireProjectRole").
						Str("project_id", projectID.String()).
						Str("user_id", userID.String()).
						Str("user_role", member.Role).
						Strs("required_roles", allowedRoles).
						Str("request_id", GetRequestID(c)).
						Msg("access denied — insufficient project role")
					return errs.NewForbiddenError("You do not have the required project role to perform this action", false)
				}
			}

			// Set project role in context for downstream handlers
			c.Set("project_role", member.Role)

			return next(c)
		}
	}
}

// containsRole checks if a role is within a list of allowed roles.
func containsRole(allowed []string, role string) bool {
	for _, r := range allowed {
		if r == role {
			return true
		}
	}
	return false
}
