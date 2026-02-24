// Package job provides background job processing using async task queues for emails and domain events.
package job

import (
	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type JobService struct {
	Client           *asynq.Client
	server           *asynq.Server
	logger           *zerolog.Logger
	eventHandlerFunc asynq.HandlerFunc
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr: redisAddr,
	}, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"default":  3, // 3 workers for default queue
			"critical": 6, // 6 workers for critical queue
			"low":      2, // 2 workers for low priority queue
		},
	})

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

// RegisterEventHandler sets the unified event handler function that will
// be registered for all domain event task types.
func (js *JobService) RegisterEventHandler(handler asynq.HandlerFunc) {
	js.eventHandlerFunc = handler
}

func (js *JobService) Start() error {
	// register task handlers

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, js.handleWelcomeEmailTask)

	// Register domain event handlers
	if js.eventHandlerFunc != nil {
		eventTypes := []string{
			"event:issue.created",
			"event:issue.updated",
			"event:issue.assigned",
			"event:issue.unassigned",
			"event:issue.status_changed",
			"event:issue.deleted",
			"event:comment.created",
			"event:comment.updated",
			"event:comment.deleted",
			"event:label.added",
			"event:label.removed",
		}
		for _, t := range eventTypes {
			mux.HandleFunc(t, js.eventHandlerFunc)
		}
	}

	js.logger.Info().Msg("Starting background job Server")
	if err := js.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (js *JobService) Stop() {
	js.logger.Info().Msg("Stopping background job Server")
	js.server.Shutdown()
	js.Client.Close()
}
