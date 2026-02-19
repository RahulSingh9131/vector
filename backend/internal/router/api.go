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

		// Project routes (nested under organizations)
		projects := orgs.Group("/:orgId/projects")
		h.Project.RegisterRoutes(projects)

		// Label routes (nested under projects)
		labels := projects.Group("/:projectId/labels")
		h.Label.RegisterRoutes(labels)

		// Issue routes (nested under projects)
		issues := projects.Group("/:projectId/issues")
		h.Issue.RegisterRoutes(issues)

		// Issue-label routes (nested under issues)
		issueLabels := issues.Group("/:issueId/labels")
		h.Label.RegisterIssueLabelRoutes(issueLabels)

		// Comment routes (nested under issues)
		comments := issues.Group("/:issueId/comments")
		h.Comment.RegisterRoutes(comments)
	}
}
