package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/labstack/echo/v4"
)

type WebhookHandler struct {
	server   *server.Server
	services *service.Services
}

func NewWebhookHandler(s *server.Server, services *service.Services) *WebhookHandler {
	return &WebhookHandler{
		server:   s,
		services: services,
	}
}

// HandleClerkWebhook processes incoming webhooks from Clerk
func (h *WebhookHandler) HandleClerkWebhook(c echo.Context) error {
	// TODO: Implement signature verification using Svix or Clerk's recommended method
	// For now, we'll log the event and process the payload
	
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
	}

	var event struct {
		Data json.RawMessage `json:"data"`
		Type string          `json:"type"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to unmarshal clerk webhook event")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook payload")
	}

	h.server.Logger.Info().Str("event_type", event.Type).Msg("received clerk webhook")

	// Handle different event types
	switch event.Type {
	case "user.created", "user.updated":
		// Handle user sync logic here if needed beyond the middleware automation
		h.server.Logger.Info().Msgf("processing %s event", event.Type)
	case "organization.created", "organization.updated":
		h.server.Logger.Info().Msgf("processing %s event", event.Type)
	}

	return c.NoContent(http.StatusOK)
}
