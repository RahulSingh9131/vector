package middleware

import "github.com/RahulSingh9131/vector/internal/server"

type RateLimitingMiddleware struct {
	server *server.Server
}

func NewRateLimitingMiddleware(s *server.Server) *RateLimitingMiddleware {
	return &RateLimitingMiddleware{
		server: s,
	}
}

func (r *RateLimitingMiddleware) RecordRateLimiting(endpoint string) {
	if r.server.LoggerService != nil && r.server.LoggerService.GetApplication() != nil {
		r.server.LoggerService.GetApplication().RecordCustomEvent("RateLimitHit", map[string]interface{}{
			"endpoint": endpoint,
		})
	}
}
