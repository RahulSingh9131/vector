// Package router defines the HTTP route registrations for the API.
package router

import (
	"net/http"

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

	// Protected routes — all require authentication
	api := r.Group("", m.Auth.RequireAuth)
	{
		// ─── User routes (no org scoping) ───
		api.GET("/me", handler.Handle(h.User.Handler, h.User.GetCurrentUser, http.StatusOK, &handler.EmptyRequest{}))
		api.PUT("/me", handler.Handle(h.User.Handler, h.User.UpdateCurrentUser, http.StatusOK, &handler.UpdateCurrentUserRequest{}))
		api.GET("/users", handler.Handle(h.User.Handler, h.User.ListUsers, http.StatusOK, &handler.EmptyRequest{}))
		api.POST("/users", handler.Handle(h.User.Handler, h.User.CreateUser, http.StatusCreated, &handler.CreateUserRequest{}))
		api.GET("/users/:id", handler.Handle(h.User.Handler, h.User.GetUser, http.StatusOK, &handler.GetUserRequest{}))
		api.DELETE("/users/:id", handler.Handle(h.User.Handler, h.User.DeleteUser, http.StatusNoContent, &handler.DeleteUserRequest{}))

		// User activity feed (my activity — scoped to current user, no org needed)
		myActivity := api.Group("/me/activity")
		myActivity.GET("", handler.Handle(h.Activity.Handler, h.Activity.ListMyActivity, http.StatusOK, &handler.ListMyActivityRequest{}))

		// ─── Organization routes ───
		orgs := api.Group("/organizations")

		// Org-level routes that don't need org membership check
		// (list my orgs = scoped to current user, create org = any authenticated user)
		orgs.GET("", handler.Handle(h.Organization.Handler, h.Organization.ListOrganizations, http.StatusOK, &handler.EmptyRequest{}))
		orgs.POST("", handler.Handle(h.Organization.Handler, h.Organization.CreateOrganization, http.StatusCreated, &handler.CreateOrganizationRequest{}))

		// Org-scoped routes — require org membership (any role)
		orgScoped := orgs.Group("/:id", m.Authorization.RequireOrgMember())
		{
			orgScoped.GET("", handler.Handle(h.Organization.Handler, h.Organization.GetOrganization, http.StatusOK, &handler.GetOrganizationRequest{}))
			orgScoped.GET("/members", handler.Handle(h.Organization.Handler, h.Organization.ListMembers, http.StatusOK, &handler.ListMembersRequest{}))
		}

		// Org member management — require org admin or owner role
		orgAdmin := orgs.Group("/:id", m.Authorization.RequireOrgMember("org:owner", "org:admin"))
		{
			orgAdmin.PATCH("", handler.Handle(h.Organization.Handler, h.Organization.UpdateOrganization, http.StatusOK, &handler.UpdateOrganizationRequest{}))
			orgAdmin.POST("/members", handler.Handle(h.Organization.Handler, h.Organization.AddMember, http.StatusCreated, &handler.AddMemberRequest{}))
			orgAdmin.PATCH("/members/:userId", handler.Handle(h.Organization.Handler, h.Organization.UpdateMemberRole, http.StatusOK, &handler.UpdateMemberRoleRequest{}))
			orgAdmin.DELETE("/members/:userId", handler.Handle(h.Organization.Handler, h.Organization.RemoveMember, http.StatusNoContent, &handler.RemoveMemberRequest{}))
		}

		// ─── Org activity (requires org membership) ───
		orgActivity := orgs.Group("/:orgId/activity", m.Authorization.RequireOrgMember())
		orgActivity.GET("", handler.Handle(h.Activity.Handler, h.Activity.ListOrgActivity, http.StatusOK, &handler.ListOrgActivityRequest{}))

		orgActivitySummary := orgs.Group("/:orgId/activity/summary", m.Authorization.RequireOrgMember())
		orgActivitySummary.GET("", handler.Handle(h.Activity.Handler, h.Activity.OrgActivitySummary, http.StatusOK, &handler.OrgActivitySummaryRequest{}))

		// ─── Project routes ───
		// Base group — requires org membership (validates org scope for all project routes)
		projects := orgs.Group("/:orgId/projects", m.Authorization.RequireOrgMember())

		// List projects — any org member (including guests)
		projects.GET("", handler.Handle(h.Project.Handler, h.Project.ListProjects, http.StatusOK, &handler.ListProjectsRequest{}))

		// Create project — owners and admins only
		projectsWrite := orgs.Group("/:orgId/projects", m.Authorization.RequireOrgMember("org:owner", "org:admin"))
		projectsWrite.POST("", handler.Handle(h.Project.Handler, h.Project.CreateProject, http.StatusCreated, &handler.CreateProjectRequest{}))

		// ─── Project-scoped routes (require project membership) ───
		projectScoped := projects.Group("/:projectId", m.Authorization.RequireProjectRole())

		// Project read — any project member
		projectScoped.GET("", handler.Handle(h.Project.Handler, h.Project.GetProject, http.StatusOK, &handler.GetProjectRequest{}))

		// Project write — project admin or member
		projectWrite := projects.Group("/:projectId", m.Authorization.RequireProjectRole("admin", "member"))
		projectWrite.PATCH("", handler.Handle(h.Project.Handler, h.Project.UpdateProject, http.StatusOK, &handler.UpdateProjectRequest{}))

		// Project delete — project admin only
		projectAdmin := projects.Group("/:projectId", m.Authorization.RequireProjectRole("admin"))
		projectAdmin.DELETE("", handler.Handle(h.Project.Handler, h.Project.DeleteProject, http.StatusNoContent, &handler.DeleteProjectRequest{}))

		// ─── Project member management ───
		// List members — any project member
		projectMembers := projectScoped.Group("/members")
		projectMembers.GET("", handler.Handle(h.Project.Handler, h.Project.ListMembers, http.StatusOK, &handler.ListProjectMembersRequest{}))

		// Add/update/remove members — project admin only
		projectMembersAdmin := projects.Group("/:projectId/members", m.Authorization.RequireProjectRole("admin"))
		projectMembersAdmin.POST("", handler.Handle(h.Project.Handler, h.Project.AddMember, http.StatusCreated, &handler.AddProjectMemberRequest{}))
		projectMembersAdmin.PATCH("/:userId", handler.Handle(h.Project.Handler, h.Project.UpdateMemberRole, http.StatusOK, &handler.UpdateProjectMemberRoleRequest{}))
		projectMembersAdmin.DELETE("/:userId", handler.Handle(h.Project.Handler, h.Project.RemoveMember, http.StatusNoContent, &handler.RemoveProjectMemberRequest{}))

		// ─── Label routes (under project) ───
		labelsRead := projectScoped.Group("/labels")
		labelsRead.GET("", handler.Handle(h.Label.Handler, h.Label.ListLabels, http.StatusOK, &handler.ListLabelsRequest{}))
		labelsRead.GET("/:labelId", handler.Handle(h.Label.Handler, h.Label.GetLabel, http.StatusOK, &handler.GetLabelRequest{}))

		labelsWrite := projects.Group("/:projectId/labels", m.Authorization.RequireProjectRole("admin", "member"))
		labelsWrite.POST("", handler.Handle(h.Label.Handler, h.Label.CreateLabel, http.StatusCreated, &handler.CreateLabelRequest{}))
		labelsWrite.PATCH("/:labelId", handler.Handle(h.Label.Handler, h.Label.UpdateLabel, http.StatusOK, &handler.UpdateLabelRequest{}))
		labelsWrite.DELETE("/:labelId", handler.Handle(h.Label.Handler, h.Label.DeleteLabel, http.StatusNoContent, &handler.DeleteLabelRequest{}))

		// ─── Issue routes (under project) ───
		issuesRead := projectScoped.Group("/issues")
		issuesRead.GET("", handler.Handle(h.Issue.Handler, h.Issue.ListIssues, http.StatusOK, &handler.ListIssuesRequest{}))
		issuesRead.GET("/:issueId", handler.Handle(h.Issue.Handler, h.Issue.GetIssue, http.StatusOK, &handler.GetIssueRequest{}))

		issuesWrite := projects.Group("/:projectId/issues", m.Authorization.RequireProjectRole("admin", "member"))
		issuesWrite.POST("", handler.Handle(h.Issue.Handler, h.Issue.CreateIssue, http.StatusCreated, &handler.CreateIssueRequest{}))
		issuesWrite.PATCH("/:issueId", handler.Handle(h.Issue.Handler, h.Issue.UpdateIssue, http.StatusOK, &handler.UpdateIssueRequest{}))
		issuesWrite.DELETE("/:issueId", handler.Handle(h.Issue.Handler, h.Issue.DeleteIssue, http.StatusNoContent, &handler.DeleteIssueRequest{}))
		issuesWrite.PATCH("/:issueId/assign", handler.Handle(h.Issue.Handler, h.Issue.AssignIssue, http.StatusOK, &handler.AssignIssueRequest{}))
		issuesWrite.PATCH("/:issueId/status", handler.Handle(h.Issue.Handler, h.Issue.UpdateIssueStatus, http.StatusOK, &handler.UpdateIssueStatusRequest{}))

		// ─── Issue-label routes ───
		issueLabelsRead := issuesRead.Group("/:issueId/labels")
		issueLabelsRead.GET("", handler.Handle(h.Label.Handler, h.Label.GetIssueLabels, http.StatusOK, &handler.GetIssueLabelsRequest{}))

		issueLabelsWrite := issuesWrite.Group("/:issueId/labels")
		issueLabelsWrite.POST("", handler.Handle(h.Label.Handler, h.Label.AddLabelToIssue, http.StatusCreated, &handler.AddLabelToIssueRequest{}))
		issueLabelsWrite.DELETE("/:labelId", handler.Handle(h.Label.Handler, h.Label.RemoveLabelFromIssue, http.StatusNoContent, &handler.RemoveLabelFromIssueRequest{}))

		// ─── Comment routes (under issues) ───
		commentsRead := issuesRead.Group("/:issueId/comments")
		commentsRead.GET("", handler.Handle(h.Comment.Handler, h.Comment.ListComments, http.StatusOK, &handler.ListCommentsRequest{}))
		commentsRead.GET("/:commentId", handler.Handle(h.Comment.Handler, h.Comment.GetComment, http.StatusOK, &handler.GetCommentRequest{}))

		commentsWrite := issuesWrite.Group("/:issueId/comments")
		commentsWrite.POST("", handler.Handle(h.Comment.Handler, h.Comment.CreateComment, http.StatusCreated, &handler.CreateCommentRequest{}))
		commentsWrite.PATCH("/:commentId", handler.Handle(h.Comment.Handler, h.Comment.UpdateComment, http.StatusOK, &handler.UpdateCommentRequest{}))
		commentsWrite.DELETE("/:commentId", handler.Handle(h.Comment.Handler, h.Comment.DeleteComment, http.StatusNoContent, &handler.DeleteCommentRequest{}))

		// ─── Issue activity (read-only, requires project membership) ───
		issueActivity := issuesRead.Group("/:issueId/activity")
		issueActivity.GET("", handler.Handle(h.Activity.Handler, h.Activity.ListIssueActivity, http.StatusOK, &handler.ListIssueActivityRequest{}))

		// ─── Project activity (read-only, requires project membership) ───
		projectActivity := projectScoped.Group("/activity")
		projectActivity.GET("", handler.Handle(h.Activity.Handler, h.Activity.ListProjectActivity, http.StatusOK, &handler.ListProjectActivityRequest{}))

		projectActivitySummary := projectScoped.Group("/activity/summary")
		projectActivitySummary.GET("", handler.Handle(h.Activity.Handler, h.Activity.ProjectActivitySummary, http.StatusOK, &handler.ProjectActivitySummaryRequest{}))

		// ─── Notification routes ───
		notifications := api.Group("/notifications")
		{
			notifications.GET("", handler.Handle(h.Notification.Handler, h.Notification.ListNotifications, http.StatusOK, &handler.ListNotificationsRequest{}))
			notifications.PATCH("/:id/read", handler.HandleNoContent(h.Notification.Handler, h.Notification.MarkAsRead, http.StatusOK, &handler.MarkAsReadRequest{}))
			notifications.POST("/read-all", handler.HandleNoContent(h.Notification.Handler, h.Notification.MarkAllAsRead, http.StatusOK, &handler.EmptyRequest{}))
			notifications.GET("/stream", h.Notification.StreamNotifications)
		}

		// ─── Search routes ───
		api.GET("/search", handler.Handle(h.Search.Handler, h.Search.Search, http.StatusOK, &handler.SearchRequest{}))
	}
}
