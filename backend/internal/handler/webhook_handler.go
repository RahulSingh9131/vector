package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RahulSingh9131/vector/internal/middleware"
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

	logger := middleware.GetLogger(c).With().
		Str("operation", "clerk_webhook").
		Logger()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
	}

	var event struct {
		Data json.RawMessage `json:"data"`
		Type string          `json:"type"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal clerk webhook event")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook payload")
	}

	logger.Info().Str("event_type", event.Type).Msg("received clerk webhook")

	// Handle different event types
	switch event.Type {
	case "user.created", "user.updated":
		// Handle user sync logic here if needed beyond the middleware automation
		logger.Info().Msgf("processing %s event", event.Type)
	case "organization.created", "organization.updated":
		logger.Info().Msgf("processing %s event", event.Type)
	default:
		logger.Warn().Str("event_type", event.Type).Msg("unhandled webhook event type")
	}

	return c.NoContent(http.StatusOK)
}
