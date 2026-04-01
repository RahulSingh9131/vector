package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RahulSingh9131/vector/internal/middleware"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	Handler
	server   *server.Server
	services *service.Services
}

func NewNotificationHandler(s *server.Server, services *service.Services) *NotificationHandler {
	return &NotificationHandler{
		Handler:  NewHandler(s),
		server:   s,
		services: services,
	}
}

// StreamNotifications handles the Server-Sent Events (SSE) stream for real-time notifications.
func (h *NotificationHandler) StreamNotifications(c echo.Context) error {
	w := c.Response().Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Subscribe to real-time updates
	notificationChan := h.server.NotificationManager.Subscribe(userID)
	defer h.server.NotificationManager.Unsubscribe(userID, notificationChan)

	// Send an initial "connected" comment to keep the connection alive
	fmt.Fprintf(w, ": connected\n\n")
	c.Response().Flush()

	ctx := c.Request().Context()

	for {
		select {
		case <-ctx.Done():
			return nil
		case n := <-notificationChan:
			// Marshal to JSON for SSE
			data, err := json.Marshal(n)
			if err != nil {
				h.server.Logger.Error().Err(err).Msg("failed to marshal notification for SSE")
				continue
			}

			// Format as SSE message: "data: {JSON}\n\n"
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			c.Response().Flush()
		}
	}
}

// ListNotificationsRequest defines the request for listing notifications.
type ListNotificationsRequest struct {
	IsRead *bool `query:"is_read"`
	Limit  int   `query:"limit"`
	Offset int   `query:"offset"`
}

func (r *ListNotificationsRequest) Validate() error {
	return nil
}

// ListNotifications returns a list of notifications for the current user.
func (h *NotificationHandler) ListNotifications(c echo.Context, req *ListNotificationsRequest) ([]models.Notification, error) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return nil, err
	}

	filters := models.NotificationFilters{
		UserID: &userID,
		IsRead: req.IsRead,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if filters.Limit == 0 {
		filters.Limit = 20
	}

	return h.services.Notification.ListNotifications(c.Request().Context(), userID, filters)
}

// MarkAsReadRequest defines the request for marking a notification as read.
type MarkAsReadRequest struct {
	ID string `param:"id"`
}

func (r *MarkAsReadRequest) Validate() error {
	if _, err := uuid.Parse(r.ID); err != nil {
		return fmt.Errorf("invalid notification ID")
	}
	return nil
}

// MarkAsRead marks a notification as read.
func (h *NotificationHandler) MarkAsRead(c echo.Context, req *MarkAsReadRequest) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	id, _ := uuid.Parse(req.ID)

	// Verify ownership
	n, err := h.services.Notification.GetNotification(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if n == nil || n.UserID != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Notification not found")
	}

	return h.services.Notification.MarkAsRead(c.Request().Context(), id)
}

// MarkAllAsRead marks all notifications for the user as read.
func (h *NotificationHandler) MarkAllAsRead(c echo.Context, _ *EmptyRequest) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	return h.services.Notification.MarkAllAsRead(c.Request().Context(), userID)
}
