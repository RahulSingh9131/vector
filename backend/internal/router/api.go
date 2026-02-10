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
		h.User.RegisterRoutes(api)

		// Organization routes
		orgs := api.Group("/organizations")
		h.Organization.RegisterRoutes(orgs)
	}
}
