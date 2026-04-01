package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/RahulSingh9131/vector/internal/database"
	"github.com/RahulSingh9131/vector/internal/events"
	"github.com/RahulSingh9131/vector/internal/handler"
	"github.com/RahulSingh9131/vector/internal/logger"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/router"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/hibiken/asynq"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize New Relic logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "local" {
		if err = database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Error().Err(err).Msg("failed to migrate database")
			return
		}
	}

	// Initialize server
	var srv *server.Server
	srv, err = server.New(cfg, &log, loggerService)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize server")
		return
	}

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		log.Error().Err(serviceErr).Msg("could not create services")
		return
	}
	handlers := handler.NewHandlers(srv, services)

	// Register event handlers
	activitySubscriber := events.NewActivitySubscriber(repos.Activity, &log)
	notificationSubscriber := events.NewNotificationSubscriber(srv, repos, &log)

	srv.Job.RegisterEventHandler(func(ctx context.Context, task *asynq.Task) error {
		// Run activity subscriber
		if err := activitySubscriber.HandleEvent(ctx, task); err != nil {
			log.Error().Err(err).Msg("activity subscriber failed")
		}
		// Run notification subscriber
		if err := notificationSubscriber.HandleEvent(ctx, task); err != nil {
			log.Error().Err(err).Msg("notification subscriber failed")
		}
		return nil
	})

	// Start job server (after event handlers are registered)
	if err := srv.Job.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start job server")
		return
	}

	// Initialize router
	r := router.NewRouter(srv, handlers, services, repos)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start server")
			stop() // Signal main thread to exit
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	log.Info().Msg("server exited properly")
}
