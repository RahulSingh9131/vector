// Package service provides business logic and orchestration for the application's core operations.
package service

import (
	"context"
	"math"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	"github.com/google/uuid"
)

// ActivityService handles business logic for activity logging
type ActivityService struct {
	server       *server.Server
	activityRepo *repository.ActivityRepository
}

// NewActivityService creates a new activity service
func NewActivityService(s *server.Server, repos *repository.Repositories) *ActivityService {
	return &ActivityService{
		server:       s,
		activityRepo: repos.Activity,
	}
}

// Record records a new activity entry. This is called synchronously from
// other services after their operations succeed. Errors are logged but
// not propagated to avoid failing the primary operation.
func (s *ActivityService) Record(ctx context.Context, params models.CreateActivityParams) {
	_, err := s.activityRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("action", params.Action).
			Str("entity_type", params.EntityType).
			Str("entity_id", params.EntityID.String()).
			Msg("failed to record activity")
	}
}

// ListByIssue retrieves paginated activities for an issue with optional filters
func (s *ActivityService) ListByIssue(ctx context.Context, issueID uuid.UUID, page, limit int, filters models.ActivityFilters) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	activities, total, err := s.activityRepo.ListByIssue(ctx, issueID, page, limit, filters)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse[models.ActivityWithActor]{
		Data:       activities,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ListByProject retrieves paginated activities for a project with optional filters
func (s *ActivityService) ListByProject(ctx context.Context, projectID uuid.UUID, page, limit int, filters models.ActivityFilters) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	activities, total, err := s.activityRepo.ListByProject(ctx, projectID, page, limit, filters)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse[models.ActivityWithActor]{
		Data:       activities,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
