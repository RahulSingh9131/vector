// Package service provides business logic and orchestration for the application's core operations.
package service

import (
	"context"
	"math"

	"github.com/RahulSingh9131/vector/internal/errs"
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

// ListByActor retrieves paginated activities performed by a specific user with optional filters
func (s *ActivityService) ListByActor(ctx context.Context, actorID uuid.UUID, page, limit int, filters models.ActivityFilters) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	activities, total, err := s.activityRepo.ListByActor(ctx, actorID, page, limit, filters)
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

// ListByOrganization retrieves paginated activities across all projects in an org that the user is a member of
func (s *ActivityService) ListByOrganization(ctx context.Context, orgID, userID uuid.UUID, page, limit int, filters models.ActivityFilters) (*models.PaginatedResponse[models.ActivityWithActor], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	activities, total, err := s.activityRepo.ListByOrganization(ctx, orgID, userID, page, limit, filters)
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

// Allowed values for group_by and interval parameters
var allowedGroupBy = map[string]bool{
	"action":      true,
	"entity_type": true,
	"actor_id":    true,
	"date":        true,
}

var allowedIntervals = map[string]bool{
	"day":   true,
	"week":  true,
	"month": true,
}

// SummaryByProject returns aggregated activity counts for a project
func (s *ActivityService) SummaryByProject(ctx context.Context, projectID uuid.UUID, groupBy, interval string, filters models.ActivityFilters) (*models.ActivitySummaryResponse, error) {
	if groupBy == "" {
		groupBy = "action"
	}
	if !allowedGroupBy[groupBy] {
		return nil, errs.NewBadRequestError(
			"Invalid group_by value. Must be one of: action, entity_type, actor_id, date",
			true, nil, nil, nil,
		)
	}

	if groupBy == "date" {
		if interval == "" {
			interval = "day"
		}
		if !allowedIntervals[interval] {
			return nil, errs.NewBadRequestError(
				"Invalid interval value. Must be one of: day, week, month",
				true, nil, nil, nil,
			)
		}
	}

	result, err := s.activityRepo.SummaryByProject(ctx, projectID, groupBy, interval, filters)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	return result, nil
}

// SummaryByOrganization returns aggregated activity counts across an org (scoped to user's projects)
func (s *ActivityService) SummaryByOrganization(ctx context.Context, orgID, userID uuid.UUID, groupBy, interval string, filters models.ActivityFilters) (*models.ActivitySummaryResponse, error) {
	if groupBy == "" {
		groupBy = "action"
	}
	if !allowedGroupBy[groupBy] {
		return nil, errs.NewBadRequestError(
			"Invalid group_by value. Must be one of: action, entity_type, actor_id, date",
			true, nil, nil, nil,
		)
	}

	if groupBy == "date" {
		if interval == "" {
			interval = "day"
		}
		if !allowedIntervals[interval] {
			return nil, errs.NewBadRequestError(
				"Invalid interval value. Must be one of: day, week, month",
				true, nil, nil, nil,
			)
		}
	}

	result, err := s.activityRepo.SummaryByOrganization(ctx, orgID, userID, groupBy, interval, filters)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	return result, nil
}
