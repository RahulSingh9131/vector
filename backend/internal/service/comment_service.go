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

// CommentService handles business logic for comments
type CommentService struct {
	server       *server.Server
	commentRepo  *repository.CommentRepository
	memberRepo   *repository.ProjectMemberRepository
	issueRepo    *repository.IssueRepository
	activityRepo *repository.ActivityRepository
}

// NewCommentService creates a new comment service
func NewCommentService(s *server.Server, repos *repository.Repositories) *CommentService {
	return &CommentService{
		server:       s,
		commentRepo:  repos.Comment,
		memberRepo:   repos.ProjectMember,
		issueRepo:    repos.Issue,
		activityRepo: repos.Activity,
	}
}

// CreateComment creates a new comment on an issue
func (s *CommentService) CreateComment(ctx context.Context, issueID, authorID uuid.UUID, params models.CreateCommentParams) (*models.CommentWithAuthor, error) {
	s.server.Logger.Debug().
		Str("issue_id", issueID.String()).
		Str("author_id", authorID.String()).
		Msg("creating comment")

	params.IssueID = issueID
	params.AuthorID = authorID

	// Validate parent comment if provided (single-level threading)
	if params.ParentCommentID != nil {
		parent, err := s.commentRepo.GetByID(ctx, *params.ParentCommentID)
		if err != nil {
			return nil, sqlerr.HandleError(err)
		}
		if parent == nil {
			return nil, errs.NewNotFoundError("Parent comment not found", true, nil)
		}

		// Parent must belong to the same issue
		if parent.IssueID != issueID {
			return nil, errs.NewBadRequestError(
				"Parent comment does not belong to this issue",
				true, nil, nil, nil,
			)
		}

		// Parent must be a top-level comment (single-level threading)
		if parent.ParentCommentID != nil {
			return nil, errs.NewBadRequestError(
				"Cannot reply to a reply. Replies are only allowed on top-level comments",
				true, nil, nil, nil,
			)
		}
	}

	comment, err := s.commentRepo.Create(ctx, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to create comment")
		return nil, sqlerr.HandleError(err)
	}

	// Fetch the full comment with author details
	commentWithAuthor, err := s.commentRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("comment_id", comment.ID.String()).
		Str("issue_id", issueID.String()).
		Msg("comment created successfully")

	// Record activity
	issue, _ := s.issueRepo.GetByID(ctx, issueID)
	if issue != nil {
		bodyPreview := params.Body
		if len(bodyPreview) > 100 {
			bodyPreview = bodyPreview[:100] + "..."
		}
		s.recordActivity(ctx, models.CreateActivityParams{
			ProjectID:  issue.ProjectID,
			IssueID:    &issueID,
			ActorID:    authorID,
			Action:     "comment.created",
			EntityType: "comment",
			EntityID:   comment.ID,
			Metadata:   map[string]interface{}{"body_preview": bodyPreview},
		})
	}

	return commentWithAuthor, nil
}

// GetByID retrieves a comment by ID with author details
func (s *CommentService) GetByID(ctx context.Context, commentID uuid.UUID) (*models.CommentWithAuthor, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if comment == nil {
		return nil, errs.NewNotFoundError("Comment not found", true, nil)
	}

	return comment, nil
}

// ListByIssue retrieves threaded comments for an issue
func (s *CommentService) ListByIssue(ctx context.Context, issueID uuid.UUID, page, limit int) (*models.PaginatedResponse[models.CommentThread], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	threads, total, err := s.commentRepo.ListByIssue(ctx, issueID, page, limit)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("issue_id", issueID.String()).
			Msg("failed to list comments")
		return nil, sqlerr.HandleError(err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse[models.CommentThread]{
		Data:       threads,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateComment updates a comment (author only)
func (s *CommentService) UpdateComment(ctx context.Context, commentID, userID uuid.UUID, params models.UpdateCommentParams) (*models.CommentWithAuthor, error) {
	s.server.Logger.Debug().
		Str("comment_id", commentID.String()).
		Str("user_id", userID.String()).
		Msg("updating comment")

	// Check authorization — only author can edit
	authorID, err := s.commentRepo.GetAuthorID(ctx, commentID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}
	if authorID == nil {
		return nil, errs.NewNotFoundError("Comment not found", true, nil)
	}
	if *authorID != userID {
		return nil, errs.NewForbiddenError("Only the comment author can edit this comment", true)
	}

	comment, err := s.commentRepo.Update(ctx, commentID, params)
	if err != nil {
		s.server.Logger.Error().Err(err).
			Str("comment_id", commentID.String()).
			Msg("failed to update comment")
		return nil, sqlerr.HandleError(err)
	}
	if comment == nil {
		return nil, errs.NewNotFoundError("Comment not found or already deleted", true, nil)
	}

	// Fetch with author details
	commentWithAuthor, err := s.commentRepo.GetByID(ctx, comment.ID)
	if err != nil {
		return nil, sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("comment_id", commentID.String()).
		Msg("comment updated successfully")

	// Record activity
	if commentWithAuthor != nil {
		issue, _ := s.issueRepo.GetByID(ctx, commentWithAuthor.IssueID)
		if issue != nil {
			s.recordActivity(ctx, models.CreateActivityParams{
				ProjectID:  issue.ProjectID,
				IssueID:    &commentWithAuthor.IssueID,
				ActorID:    userID,
				Action:     "comment.updated",
				EntityType: "comment",
				EntityID:   commentID,
			})
		}
	}

	return commentWithAuthor, nil
}

// DeleteComment soft-deletes a comment (author or project admin)
func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID, projectID uuid.UUID, issueID uuid.UUID) error {
	s.server.Logger.Debug().
		Str("comment_id", commentID.String()).
		Str("user_id", userID.String()).
		Msg("deleting comment")

	// Check authorization — author or project admin can delete
	authorID, err := s.commentRepo.GetAuthorID(ctx, commentID)
	if err != nil {
		return sqlerr.HandleError(err)
	}
	if authorID == nil {
		return errs.NewNotFoundError("Comment not found", true, nil)
	}

	isAuthor := *authorID == userID

	if !isAuthor {
		// Check if user is a project admin
		member, memberErr := s.memberRepo.GetMember(ctx, projectID, userID)
		if memberErr != nil {
			return sqlerr.HandleError(memberErr)
		}
		if member == nil || member.Role != "admin" {
			return errs.NewForbiddenError("Only the comment author or a project admin can delete this comment", true)
		}
	}

	if err := s.commentRepo.SoftDelete(ctx, commentID); err != nil {
		s.server.Logger.Error().Err(err).
			Str("comment_id", commentID.String()).
			Msg("failed to delete comment")
		return sqlerr.HandleError(err)
	}

	s.server.Logger.Info().
		Str("comment_id", commentID.String()).
		Msg("comment soft-deleted successfully")

	// Record activity
	s.recordActivity(ctx, models.CreateActivityParams{
		ProjectID:  projectID,
		IssueID:    &issueID,
		ActorID:    userID,
		Action:     "comment.deleted",
		EntityType: "comment",
		EntityID:   commentID,
	})

	return nil
}

// recordActivity is a helper that logs errors but doesn't propagate them
func (s *CommentService) recordActivity(ctx context.Context, params models.CreateActivityParams) {
	if _, err := s.activityRepo.Create(ctx, params); err != nil {
		s.server.Logger.Error().Err(err).
			Str("action", params.Action).
			Msg("failed to record activity")
	}
}
