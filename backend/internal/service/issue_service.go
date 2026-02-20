package service

import (
	"context"
	"fmt"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	"github.com/google/uuid"
)

// Allowed issue statuses
var allowedIssueStatuses = map[string]bool{
	"backlog":     true,
	"todo":        true,
	"in_progress": true,
	"in_review":   true,
	"in_dev":      true,
	"in_prod":     true,
	"cancelled":   true,
}

// Allowed issue priorities
var allowedIssuePriorities = map[string]bool{
	"urgent": true,
	"high":   true,
	"medium": true,
	"low":    true,
	"none":   true,
}

// Allowed issue types
var allowedIssueTypes = map[string]bool{
	"task":  true,
	"bug":   true,
	"story": true,
	"epic":  true,
}

// IssueService handles business logic for issues
type IssueService struct {
	server       *server.Server
	issueRepo    *repository.IssueRepository
	projectRepo  *repository.ProjectRepository
	memberRepo   *repository.ProjectMemberRepository
	activityRepo *repository.ActivityRepository
}

// NewIssueService creates a new issue service
func NewIssueService(s *server.Server, repos *repository.Repositories) *IssueService {
	return &IssueService{
		server:       s,
		issueRepo:    repos.Issue,
		projectRepo:  repos.Project,
		memberRepo:   repos.ProjectMember,
		activityRepo: repos.Activity,
	}
}

// CreateIssue creates a new issue within a project
func (s *IssueService) CreateIssue(ctx context.Context, projectID, reporterID uuid.UUID, params models.CreateIssueParams) (*models.Issue, error) {
	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Str("reporter_id", reporterID.String()).
		Str("title", params.Title).
		Msg("creating issue")

	// Set project and reporter
	params.ProjectID = projectID
	params.ReporterID = reporterID

	// Validate priority
	if params.Priority == "" {
		params.Priority = "none"
	}
	if !allowedIssuePriorities[params.Priority] {
		return nil, errs.NewBadRequestError(
			"Invalid priority. Must be one of: urgent, high, medium, low, none",
			true, nil,
			[]errs.FieldError{{Field: "priority", Error: "must be one of: urgent, high, medium, low, none"}},
			nil,
		)
	}

	// Validate type
	if params.Type == "" {
		params.Type = "task"
	}
	if !allowedIssueTypes[params.Type] {
		return nil, errs.NewBadRequestError(
			"Invalid type. Must be one of: task, bug, story, epic",
			true, nil,
			[]errs.FieldError{{Field: "type", Error: "must be one of: task, bug, story, epic"}},
			nil,
		)
	}

	// Validate assignee is a project member (if provided)
	if params.AssigneeID != nil {
		member, err := s.memberRepo.GetMember(ctx, projectID, *params.AssigneeID)
		if err != nil {
			return nil, sqlerr.HandleError(err)
		}
		if member == nil {
			return nil, errs.NewBadRequestError(
				"Assignee is not a member of this project",
				true, nil,
				[]errs.FieldError{{Field: "assignee_id", Error: "must be a project member"}},
				nil,
			)
		}
	}

	// Get project to generate issue key
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if project == nil {
		return nil, errs.NewNotFoundError("Project not found", false, nil)
	}

	// Atomically increment the issue counter and generate the key
	counter, err := s.projectRepo.IncrementIssueCounter(ctx, projectID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	issueKey := fmt.Sprintf("%s-%d", project.Identifier, counter)

	issue, err := s.issueRepo.Create(ctx, params, issueKey)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Str("title", params.Title).
			Msg("failed to create issue")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issue.ID.String()).
		Str("issue_key", issue.IssueKey).
		Str("project_id", projectID.String()).
		Msg("issue created successfully")

	// Record activity
	s.recordActivity(ctx, models.CreateActivityParams{
		ProjectID:  projectID,
		IssueID:    &issue.ID,
		ActorID:    reporterID,
		Action:     "issue.created",
		EntityType: "issue",
		EntityID:   issue.ID,
		NewValue: map[string]interface{}{
			"title":    issue.Title,
			"status":   issue.Status,
			"priority": issue.Priority,
			"type":     issue.Type,
		},
	})

	return issue, nil
}

// GetByID retrieves an issue by ID with details
func (s *IssueService) GetByID(ctx context.Context, id uuid.UUID) (*models.IssueWithDetails, error) {
	issue, err := s.issueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if issue == nil {
		return nil, errs.NewNotFoundError("Issue not found", false, nil)
	}

	return issue, nil
}

// ListIssues retrieves issues for a project with optional filters
func (s *IssueService) ListIssues(ctx context.Context, projectID uuid.UUID, filters models.IssueFilters) (*models.PaginatedResponse[models.IssueWithDetails], error) {
	// Validate filter values if provided
	if filters.Status != nil && !allowedIssueStatuses[*filters.Status] {
		return nil, errs.NewBadRequestError(
			"Invalid status filter",
			true, nil,
			[]errs.FieldError{{Field: "status", Error: "must be a valid issue status"}},
			nil,
		)
	}
	if filters.Priority != nil && !allowedIssuePriorities[*filters.Priority] {
		return nil, errs.NewBadRequestError(
			"Invalid priority filter",
			true, nil,
			[]errs.FieldError{{Field: "priority", Error: "must be a valid issue priority"}},
			nil,
		)
	}
	if filters.Type != nil && !allowedIssueTypes[*filters.Type] {
		return nil, errs.NewBadRequestError(
			"Invalid type filter",
			true, nil,
			[]errs.FieldError{{Field: "type", Error: "must be a valid issue type"}},
			nil,
		)
	}

	result, err := s.issueRepo.ListByProject(ctx, projectID, filters)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Msg("failed to list issues")
		return nil, sqlerr.HandleError(err)
	}

	return result, nil
}

// UpdateIssue updates an issue
func (s *IssueService) UpdateIssue(ctx context.Context, issueID, actorID uuid.UUID, params models.UpdateIssueParams) (*models.Issue, error) {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Msg("updating issue")

	// Validate status if provided
	if params.Status != nil && !allowedIssueStatuses[*params.Status] {
		return nil, errs.NewBadRequestError(
			"Invalid status",
			true, nil,
			[]errs.FieldError{{Field: "status", Error: "must be a valid issue status"}},
			nil,
		)
	}

	// Validate priority if provided
	if params.Priority != nil && !allowedIssuePriorities[*params.Priority] {
		return nil, errs.NewBadRequestError(
			"Invalid priority",
			true, nil,
			[]errs.FieldError{{Field: "priority", Error: "must be a valid issue priority"}},
			nil,
		)
	}

	// Validate type if provided
	if params.Type != nil && !allowedIssueTypes[*params.Type] {
		return nil, errs.NewBadRequestError(
			"Invalid type",
			true, nil,
			[]errs.FieldError{{Field: "type", Error: "must be a valid issue type"}},
			nil,
		)
	}

	// Fetch existing issue for activity tracking and validation
	existingIssue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingIssue == nil {
		return nil, errs.NewNotFoundError("Issue not found", false, nil)
	}

	// Validate assignee is a project member (if changing assignee)
	if params.AssigneeID != nil {
		member, memberErr := s.memberRepo.GetMember(ctx, existingIssue.ProjectID, *params.AssigneeID)
		if memberErr != nil {
			return nil, sqlerr.HandleError(memberErr)
		}
		if member == nil {
			return nil, errs.NewBadRequestError(
				"Assignee is not a member of this project",
				true, nil,
				[]errs.FieldError{{Field: "assignee_id", Error: "must be a project member"}},
				nil,
			)
		}
	}

	// Build old/new values for activity tracking
	oldValue := map[string]interface{}{}
	newValue := map[string]interface{}{}
	if params.Title != nil {
		oldValue["title"] = existingIssue.Title
		newValue["title"] = *params.Title
	}
	if params.Status != nil {
		oldValue["status"] = existingIssue.Status
		newValue["status"] = *params.Status
	}
	if params.Priority != nil {
		oldValue["priority"] = existingIssue.Priority
		newValue["priority"] = *params.Priority
	}
	if params.Type != nil {
		oldValue["type"] = existingIssue.Type
		newValue["type"] = *params.Type
	}

	issue, err := s.issueRepo.Update(ctx, issueID, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to update issue")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Msg("issue updated successfully")

	// Record activity
	if len(oldValue) > 0 {
		s.recordActivity(ctx, models.CreateActivityParams{
			ProjectID:  existingIssue.ProjectID,
			IssueID:    &issueID,
			ActorID:    actorID,
			Action:     "issue.updated",
			EntityType: "issue",
			EntityID:   issueID,
			OldValue:   oldValue,
			NewValue:   newValue,
		})
	}

	return issue, nil
}

// AssignIssue assigns or unassigns an issue
func (s *IssueService) AssignIssue(ctx context.Context, issueID, actorID uuid.UUID, assigneeID *uuid.UUID) (*models.Issue, error) {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Msg("assigning issue")

	// Get the issue to find its project
	existingIssue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingIssue == nil {
		return nil, errs.NewNotFoundError("Issue not found", false, nil)
	}

	// Validate assignee is a project member (if assigning, not unassigning)
	if assigneeID != nil {
		member, memberErr := s.memberRepo.GetMember(ctx, existingIssue.ProjectID, *assigneeID)
		if memberErr != nil {
			return nil, sqlerr.HandleError(memberErr)
		}
		if member == nil {
			return nil, errs.NewBadRequestError(
				"Assignee is not a member of this project",
				true, nil,
				[]errs.FieldError{{Field: "assignee_id", Error: "must be a project member"}},
				nil,
			)
		}
	}

	issue, err := s.issueRepo.Update(ctx, issueID, models.UpdateIssueParams{
		AssigneeID: assigneeID,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to assign issue")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Msg("issue assigned successfully")

	// Record activity
	action := "issue.assigned"
	if assigneeID == nil {
		action = "issue.unassigned"
	}
	oldAssignee := map[string]interface{}{"assignee_id": existingIssue.AssigneeID}
	newAssignee := map[string]interface{}{"assignee_id": assigneeID}
	s.recordActivity(ctx, models.CreateActivityParams{
		ProjectID:  existingIssue.ProjectID,
		IssueID:    &issueID,
		ActorID:    actorID,
		Action:     action,
		EntityType: "issue",
		EntityID:   issueID,
		OldValue:   oldAssignee,
		NewValue:   newAssignee,
	})

	return issue, nil
}

// UpdateStatus updates an issue's status
func (s *IssueService) UpdateStatus(ctx context.Context, issueID, actorID uuid.UUID, status string) (*models.Issue, error) {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Str("status", status).
		Msg("updating issue status")

	if !allowedIssueStatuses[status] {
		return nil, errs.NewBadRequestError(
			"Invalid status. Must be one of: backlog, todo, in_progress, in_review, in_dev, in_prod, cancelled",
			true, nil,
			[]errs.FieldError{{Field: "status", Error: "must be one of: backlog, todo, in_progress, in_review, in_dev, in_prod, cancelled"}},
			nil,
		)
	}

	// Fetch existing issue for activity tracking
	existingIssue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if existingIssue == nil {
		return nil, errs.NewNotFoundError("Issue not found", false, nil)
	}

	issue, err := s.issueRepo.Update(ctx, issueID, models.UpdateIssueParams{
		Status: &status,
	})
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to update issue status")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Str("status", status).
		Msg("issue status updated successfully")

	// Record activity
	s.recordActivity(ctx, models.CreateActivityParams{
		ProjectID:  existingIssue.ProjectID,
		IssueID:    &issueID,
		ActorID:    actorID,
		Action:     "issue.status_changed",
		EntityType: "issue",
		EntityID:   issueID,
		OldValue:   map[string]interface{}{"status": existingIssue.Status},
		NewValue:   map[string]interface{}{"status": status},
	})

	return issue, nil
}

// DeleteIssue deletes an issue
func (s *IssueService) DeleteIssue(ctx context.Context, issueID, actorID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Msg("deleting issue")

	// Fetch issue for activity metadata before deleting
	existingIssue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if existingIssue == nil {
		return errs.NewNotFoundError("Issue not found", false, nil)
	}

	if err := s.issueRepo.Delete(ctx, issueID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to delete issue")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Msg("issue deleted successfully")

	// Record activity (issue_id is null since the issue is deleted)
	s.recordActivity(ctx, models.CreateActivityParams{
		ProjectID:  existingIssue.ProjectID,
		ActorID:    actorID,
		Action:     "issue.deleted",
		EntityType: "issue",
		EntityID:   issueID,
		Metadata:   map[string]interface{}{"title": existingIssue.Title, "issue_key": existingIssue.IssueKey},
	})

	return nil
}

// recordActivity is a helper that logs errors but doesn't propagate them
func (s *IssueService) recordActivity(ctx context.Context, params models.CreateActivityParams) {
	if _, err := s.activityRepo.Create(ctx, params); err != nil {
		s.server.Logger.Error().Err(err).
			Str("action", params.Action).
			Msg("failed to record activity")
	}
}
