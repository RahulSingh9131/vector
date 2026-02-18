package service

import (
	"context"
	"regexp"

	"github.com/RahulSingh9131/vector/internal/errs"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/repository"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/RahulSingh9131/vector/internal/sqlerr"
	"github.com/google/uuid"
)

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

// LabelService handles business logic for labels
type LabelService struct {
	server    *server.Server
	labelRepo *repository.LabelRepository
	issueRepo *repository.IssueRepository
}

// NewLabelService creates a new label service
func NewLabelService(s *server.Server, repos *repository.Repositories) *LabelService {
	return &LabelService{
		server:    s,
		labelRepo: repos.Label,
		issueRepo: repos.Issue,
	}
}

// CreateLabel creates a new label within a project
func (s *LabelService) CreateLabel(ctx context.Context, projectID uuid.UUID, params models.CreateLabelParams) (*models.Label, error) {
	s.server.Logger.Debug().
		Str("project_id", projectID.String()).
		Str("name", params.Name).
		Msg("creating label")

	// Validate hex color
	if !hexColorRegex.MatchString(params.Color) {
		return nil, errs.NewBadRequestError(
			"Invalid color format. Must be a hex color like #FF0000",
			true, nil,
			[]errs.FieldError{{Field: "color", Error: "must be a valid hex color (e.g. #FF0000)"}},
			nil,
		)
	}

	params.ProjectID = projectID

	label, err := s.labelRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Str("name", params.Name).
			Msg("failed to create label")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("label_id", label.ID.String()).
		Str("project_id", projectID.String()).
		Str("name", params.Name).
		Msg("label created successfully")

	return label, nil
}

// GetByID retrieves a label by ID
func (s *LabelService) GetByID(ctx context.Context, id uuid.UUID) (*models.Label, error) {
	label, err := s.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if label == nil {
		return nil, errs.NewNotFoundError("Label not found", false, nil)
	}

	return label, nil
}

// ListByProject retrieves all labels for a project
func (s *LabelService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Label, error) {
	labels, err := s.labelRepo.ListByProject(ctx, projectID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("project_id", projectID.String()).
			Msg("failed to list labels")
		return nil, sqlerr.HandleError(err)
	}

	return labels, nil
}

// UpdateLabel updates a label
func (s *LabelService) UpdateLabel(ctx context.Context, labelID uuid.UUID, params models.UpdateLabelParams) (*models.Label, error) {
	s.server.Logger.Debug().
		Str("label_id", labelID.String()).
		Msg("updating label")

	// Validate hex color if provided
	if params.Color != nil && !hexColorRegex.MatchString(*params.Color) {
		return nil, errs.NewBadRequestError(
			"Invalid color format. Must be a hex color like #FF0000",
			true, nil,
			[]errs.FieldError{{Field: "color", Error: "must be a valid hex color (e.g. #FF0000)"}},
			nil,
		)
	}

	label, err := s.labelRepo.Update(ctx, labelID, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("label_id", labelID.String()).
			Msg("failed to update label")
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("label_id", labelID.String()).
		Msg("label updated successfully")

	return label, nil
}

// DeleteLabel deletes a label
func (s *LabelService) DeleteLabel(ctx context.Context, labelID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("label_id", labelID.String()).
		Msg("deleting label")

	if err := s.labelRepo.Delete(ctx, labelID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("label_id", labelID.String()).
			Msg("failed to delete label")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("label_id", labelID.String()).
		Msg("label deleted successfully")

	return nil
}

// AddLabelToIssue attaches a label to an issue
func (s *LabelService) AddLabelToIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Str("label_id", labelID.String()).
		Msg("adding label to issue")

	// Verify the label exists
	label, err := s.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if label == nil {
		return errs.NewNotFoundError("Label not found", false, nil)
	}

	// Verify the issue exists and belongs to the same project
	issue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if issue == nil {
		return errs.NewNotFoundError("Issue not found", false, nil)
	}

	if issue.ProjectID != label.ProjectID {
		return errs.NewBadRequestError(
			"Label and issue must belong to the same project",
			true, nil, nil, nil,
		)
	}

	if err := s.labelRepo.AddLabelToIssue(ctx, issueID, labelID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Str("label_id", labelID.String()).
			Msg("failed to add label to issue")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Str("label_id", labelID.String()).
		Msg("label added to issue successfully")

	return nil
}

// RemoveLabelFromIssue removes a label from an issue
func (s *LabelService) RemoveLabelFromIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Str("label_id", labelID.String()).
		Msg("removing label from issue")

	if err := s.labelRepo.RemoveLabelFromIssue(ctx, issueID, labelID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Str("label_id", labelID.String()).
			Msg("failed to remove label from issue")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("issue_id", issueID.String()).
		Str("label_id", labelID.String()).
		Msg("label removed from issue successfully")

	return nil
}

// GetLabelsByIssue retrieves all labels attached to an issue
func (s *LabelService) GetLabelsByIssue(ctx context.Context, issueID uuid.UUID) ([]models.Label, error) {
	labels, err := s.labelRepo.GetLabelsByIssue(ctx, issueID)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to get labels for issue")
		return nil, sqlerr.HandleError(err)
	}

	return labels, nil
}
