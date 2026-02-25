// Package job provides background job processing using async task queues for emails and domain events.
package job

import (
	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/RahulSingh9131/vector/internal/database"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type JobService struct {
	Client           *asynq.Client
	server           *asynq.Server
	scheduler        *asynq.Scheduler
	logger           *zerolog.Logger
	db               database.DBTX
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

	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{
		Addr: redisAddr,
	}, nil)

	return &JobService{
		Client:    client,
		server:    server,
		scheduler: scheduler,
		logger:    logger,
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
			"event:project.created",
			"event:project.updated",
			"event:project.deleted",
			"event:member.added",
			"event:member.removed",
			"event:member.role_changed",
		}
		for _, t := range eventTypes {
			mux.HandleFunc(t, js.eventHandlerFunc)
		}
	}

	// Register cleanup task handler
	mux.HandleFunc(TaskActivityCleanup, js.handleActivityCleanupTask)

	// Schedule periodic activity cleanup — daily at 3 AM
	_, err := js.scheduler.Register("0 3 * * *", asynq.NewTask(TaskActivityCleanup, nil))
	if err != nil {
		return err
	}

	// Start the scheduler
	if err := js.scheduler.Start(); err != nil {
		return err
	}

	js.logger.Info().Msg("Starting background job Server")
	if err := js.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (js *JobService) Stop() {
	js.logger.Info().Msg("Stopping background job Server")
	js.scheduler.Shutdown()
	js.server.Shutdown()
	js.Client.Close()
}
