package events

import "context"

// EventBus abstracts event delivery. The default implementation uses Asynq
// (Redis-backed task queue). Future implementations can swap to NATS, Kafka,
// or any other message broker without changing service code.
type EventBus interface {
	// Publish enqueues a domain event for asynchronous processing.
	// Implementations should be fire-and-forget from the caller's perspective;
	// errors are logged but do not fail the primary operation.
	Publish(ctx context.Context, event DomainEvent) error
}
