package job

import (
	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type JobService struct {
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
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

func (js *JobService) Start() error {
	// regiter task handlers

	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, js.handleWelcomeEmailTask)

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
