package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/RahulSingh9131/vector/internal/database"
	"github.com/RahulSingh9131/vector/internal/lib/job"
	"github.com/RahulSingh9131/vector/internal/logger"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *logger.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.JobService
}

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *logger.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Redis client with New reilc integration
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	// Add New relic Redis hooks if available
	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}

	// test Redis connnection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	logger.Info().Msg("connected to the redis")

	// Start job service
	jobService := job.NewJobService(logger, cfg)
	jobService.InitHandlers(cfg, logger)

	// Start job server
	if err := jobService.Start(); err != nil {
		return nil, fmt.Errorf("failed to start job server: %w", err)
	}

	return &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
	}, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         "" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("http server is not initialized")
	}
	s.Logger.Info().Str("port", s.Config.Server.Port).Str("env", s.Config.Primary.Env).Msg("Starting HTTP server")

	return s.httpServer.ListenAndServe()
}


func (s * Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	if s.Job != nil {
		s.Job.Stop()
	}

	return nil
}