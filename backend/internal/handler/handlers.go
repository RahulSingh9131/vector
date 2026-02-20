// Package handler provides HTTP request handlers for the Vector API.
package handler

import (
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
)

type Handlers struct {
	Health       *HealthHandler
	OpenAPI      *OpenAPIHandler
	User         *UserHandler
	Organization *OrganizationHandler
	Webhook      *WebhookHandler
	Project      *ProjectHandler
	Issue        *IssueHandler
	Label        *LabelHandler
	Comment      *CommentHandler
	Activity     *ActivityHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:       NewHealthHandler(s),
		OpenAPI:      NewOpenAPIHandler(s),
		User:         NewUserHandler(s, services),
		Organization: NewOrganizationHandler(s, services),
		Webhook:      NewWebhookHandler(s, services),
		Project:      NewProjectHandler(s, services),
		Issue:        NewIssueHandler(s, services),
		Label:        NewLabelHandler(s, services),
		Comment:      NewCommentHandler(s, services),
		Activity:     NewActivityHandler(s, services),
	}
}



