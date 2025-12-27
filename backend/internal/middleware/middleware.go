package middleware

import (
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middleware struct {
	Global          *GlobalMiddlewares
	Auth            *AuthMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	RateLimiting    *RateLimitingMiddleware
}

func NewMiddlewares(s *server.Server) *Middleware {
	// get New relic application instance from server
	var nrApp *newrelic.Application
	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	return &Middleware{
		Global:          NewGlobalMiddlewares(s),
		Auth:            NewAuthMiddleware(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracingMiddleware(s, nrApp),
		RateLimiting:    NewRateLimitingMiddleware(s),
	}
}
