// Package middleware provides HTTP middleware functions for the Vector API,
// including authentication, authorization, RBAC, and request context handling.
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RahulSingh9131/vector/internal/errs"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	server   *server.Server
	services *service.Services
}

func NewAuthMiddleware(s *server.Server, services *service.Services) *AuthMiddleware {
	return &AuthMiddleware{
		server:   s,
		services: services,
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(
		clerkhttp.WithHeaderAuthorization(
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				response := map[string]string{
					"code":     "UNAUTHORIZED",
					"message":  "Unauthorized",
					"override": "false",
					"status":   "401",
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					auth.server.Logger.Error().Err(err).Str("function", "RequireAuth").Dur(
						"duration", time.Since(start)).Msg("failed to write JSON response")
				} else {
					auth.server.Logger.Error().Str("function", "RequireAuth").Dur("duration", time.Since(start)).Msg(
						"clerk authentication failed: invalid or missing token")
				}
			}))))(func(c echo.Context) error {
		start := time.Now()
		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context())

		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("clerk claims missing from context after successful auth wrapper")
			return errs.NewUnauthorizedError("Unauthorized: claims missing", false)
		}

		// Get or create user in our database
		user, err := auth.services.User.GetOrCreateFromClerk(c.Request().Context(), claims.Subject)
		if err != nil {
			auth.server.Logger.Error().
				Err(err).
				Str("function", "RequireAuth").
				Str("clerk_user_id", claims.Subject).
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("failed to get or create user")
			return errs.NewUnauthorizedError("Failed to authenticate user", false)
		}

		// Record login asynchronously (don't block request)
		go func() {
			// Use context.Background() to ensure the background task completes
			// even after the HTTP request context is canceled.
			bgCtx := context.Background()

			auth.server.Logger.Debug().
				Str("user_id", user.ID.String()).
				Msg("starting background record login task")

			if err := auth.services.User.RecordLogin(bgCtx, user.ID); err != nil {
				auth.server.Logger.Warn().
					Err(err).
					Str("user_id", user.ID.String()).
					Str("ctx_error", fmt.Sprintf("%v", bgCtx.Err())).
					Msg("failed to record login background task")
			}
		}()

		// Attach user and claims to context
		c.Set("user", user)                                               // Full user object
		c.Set("user_id", user.ID)                                         // Internal UUID
		c.Set("clerk_user_id", claims.Subject)                            // Clerk user ID
		c.Set("user_role", claims.ActiveOrganizationRole)                 // Clerk org role
		c.Set("org_id", claims.ActiveOrganizationID)                      // Clerk org ID
		c.Set("permissions", claims.Claims.ActiveOrganizationPermissions) // Clerk permissions

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", user.ID.String()).
			Str("clerk_user_id", claims.Subject).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)
	})
}
