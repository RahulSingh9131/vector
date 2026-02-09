package router

import (
	"github.com/RahulSingh9131/vector/internal/handler"
	"github.com/RahulSingh9131/vector/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerAPIRoutes(r *echo.Group, h *handler.Handlers, m *middleware.Middleware) {
	// Webhooks (typically unauthenticated or has custom signature verification)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/clerk", h.Webhook.HandleClerkWebhook)
	}

	// Protected routes
	api := r.Group("", m.Auth.RequireAuth)
	{
		// User routes
		api.GET("/me", h.User.GetCurrentUser)
		api.PUT("/me", h.User.UpdateCurrentUser)

		// Organization routes
		orgs := api.Group("/organizations")
		{
			orgs.GET("", h.Organization.ListOrganizations)
			orgs.GET("/:id", h.Organization.GetOrganization)
			orgs.GET("/:id/members", h.Organization.ListMembers)
		}
	}
}
