package middleware

import (
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/service"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middleware struct {
	Global          *GlobalMiddlewares
	Auth            *AuthMiddleware
	Authorization   *AuthorizationMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	RateLimiting    *RateLimitingMiddleware
}

func NewMiddlewares(s *server.Server, services *service.Services, repos *repository.Repositories) *Middleware {
	// get New relic application instance from server
	var nrApp *newrelic.Application
	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	return &Middleware{
		Global:          NewGlobalMiddlewares(s),
		Auth:            NewAuthMiddleware(s, services),
		Authorization:   NewAuthorizationMiddleware(s, repos),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracingMiddleware(s, nrApp),
		RateLimiting:    NewRateLimitingMiddleware(s),
	}
}
