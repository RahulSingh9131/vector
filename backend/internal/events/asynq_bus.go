package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	// TaskPrefix is prepended to event types when creating Asynq tasks.
	// e.g. "issue.created" becomes "event:issue.created"
	TaskPrefix = "event:"
)

// AsynqEventBus implements EventBus using the Asynq distributed task queue.
// Events are persisted in Redis and processed by registered task handlers
// with automatic retries and dead-letter handling.
type AsynqEventBus struct {
	client *asynq.Client
	logger *zerolog.Logger
}

// NewAsynqEventBus creates a new Asynq-backed event bus.
func NewAsynqEventBus(client *asynq.Client, logger *zerolog.Logger) *AsynqEventBus {
	return &AsynqEventBus{
		client: client,
		logger: logger,
	}
}

// Publish enqueues a domain event as an Asynq task.
func (b *AsynqEventBus) Publish(ctx context.Context, event DomainEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	taskType := TaskPrefix + event.Type
	task := asynq.NewTask(taskType, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second),
	)

	info, err := b.client.Enqueue(task)
	if err != nil {
		b.logger.Error().
			Err(err).
			Str("event_type", event.Type).
			Str("project_id", event.ProjectID.String()).
			Msg("failed to enqueue domain event")
		return fmt.Errorf("failed to enqueue event %s: %w", event.Type, err)
	}

	b.logger.Debug().
		Str("task_id", info.ID).
		Str("event_type", event.Type).
		Str("queue", info.Queue).
		Msg("domain event enqueued")

	return nil
}
