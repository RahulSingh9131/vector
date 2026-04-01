// Package realtime provides utilities for real-time communication,
// specifically focusing on notification delivery via SSE and Redis.
package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	// Redis channel for notification broadcasting across nodes
	NotificationChannel = "vector:notifications"
)

// NotificationManager handles real-time notification delivery via SSE and Redis Pub/Sub.
type NotificationManager struct {
	logger *zerolog.Logger
	redis  *redis.Client
	// Map of userID to a slice of notification channels (one for each open SSE connection)
	clients    map[uuid.UUID][]chan *models.Notification
	clientsMux sync.RWMutex
}

// NewNotificationManager creates a new notification manager.
func NewNotificationManager(logger *zerolog.Logger, r *redis.Client) *NotificationManager {
	return &NotificationManager{
		logger:  logger,
		redis:   r,
		clients: make(map[uuid.UUID][]chan *models.Notification),
	}
}

// Subscribe adds a new client channel for a user.
func (m *NotificationManager) Subscribe(userID uuid.UUID) chan *models.Notification {
	m.clientsMux.Lock()
	defer m.clientsMux.Unlock()

	ch := make(chan *models.Notification, 10)
	m.clients[userID] = append(m.clients[userID], ch)

	m.logger.Debug().
		Str("user_id", userID.String()).
		Int("active_connections", len(m.clients[userID])).
		Msg("user subscribed to real-time notifications")

	return ch
}

// Unsubscribe removes a client channel for a user.
func (m *NotificationManager) Unsubscribe(userID uuid.UUID, ch chan *models.Notification) {
	m.clientsMux.Lock()
	defer m.clientsMux.Unlock()

	channels := m.clients[userID]
	for i, c := range channels {
		if c == ch {
			m.clients[userID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}

	if len(m.clients[userID]) == 0 {
		delete(m.clients, userID)
	}

	m.logger.Debug().
		Str("user_id", userID.String()).
		Msg("user unsubscribed from real-time notifications")
}

// Publish sends a notification to a specific user.
// It persists the notification to Redis Pub/Sub for cross-node delivery.
func (m *NotificationManager) Publish(ctx context.Context, notification *models.Notification) error {
	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Publish to Redis for cross-node delivery
	if err := m.redis.Publish(ctx, NotificationChannel, payload).Err(); err != nil {
		m.logger.Error().Err(err).Msg("failed to publish notification to redis")
		return err
	}

	return nil
}

// StartRedisListener listens for notifications on the Redis Pub/Sub channel
// and delivers them to local clients.
func (m *NotificationManager) StartRedisListener(ctx context.Context) {
	pubsub := m.redis.Subscribe(ctx, NotificationChannel)
	defer pubsub.Close()

	m.logger.Info().Msg("started redis notification listener")

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			var n models.Notification
			if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
				m.logger.Error().Err(err).Msg("failed to unmarshal redis notification")
				continue
			}

			m.deliverLocal(n.UserID, &n)
		}
	}
}

// deliverLocal sends a notification to all active SSE connections for a user on THIS node.
func (m *NotificationManager) deliverLocal(userID uuid.UUID, n *models.Notification) {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	channels, ok := m.clients[userID]
	if !ok {
		return
	}

	for _, ch := range channels {
		select {
		case ch <- n:
		default:
			m.logger.Warn().
				Str("user_id", userID.String()).
				Msg("notification channel full, dropping real-time message")
		}
	}
}
