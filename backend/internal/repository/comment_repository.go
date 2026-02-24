package repository

import (
	"context"
	"errors"

	"github.com/RahulSingh9131/vector/internal/database"
	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CommentRepository handles database operations for comments
type CommentRepository struct {
	db     database.DBTX
	server *server.Server
}

// NewCommentRepository creates a new comment repository
func NewCommentRepository(s *server.Server) *CommentRepository {
	return &CommentRepository{
		db:     s.DB.Pool,
		server: s,
	}
}

// Create creates a new comment
func (r *CommentRepository) Create(ctx context.Context, params models.CreateCommentParams) (*models.Comment, error) {
	query := `
		INSERT INTO comments (issue_id, author_id, body, parent_comment_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, issue_id, author_id, body, parent_comment_id,
		          is_edited, edited_at, is_deleted, deleted_at, created_at, updated_at
	`

	var comment models.Comment
	err := r.db.QueryRow(
		ctx, query,
		params.IssueID, params.AuthorID, params.Body, params.ParentCommentID,
	).Scan(
		&comment.ID, &comment.IssueID, &comment.AuthorID, &comment.Body,
		&comment.ParentCommentID, &comment.IsEdited, &comment.EditedAt,
		&comment.IsDeleted, &comment.DeletedAt, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// GetByID retrieves a comment with author details
func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CommentWithAuthor, error) {
	query := `
		SELECT c.id, c.issue_id, c.author_id, c.body, c.parent_comment_id,
		       c.is_edited, c.edited_at, c.is_deleted, c.deleted_at, c.created_at, c.updated_at,
		       u.first_name AS author_first_name, u.last_name AS author_last_name,
		       u.avatar_url AS author_avatar_url, u.email AS author_email
		FROM comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.id = $1
	`

	var cwa models.CommentWithAuthor
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cwa.ID, &cwa.IssueID, &cwa.AuthorID, &cwa.Body,
		&cwa.ParentCommentID, &cwa.IsEdited, &cwa.EditedAt,
		&cwa.IsDeleted, &cwa.DeletedAt, &cwa.CreatedAt, &cwa.UpdatedAt,
		&cwa.AuthorFirstName, &cwa.AuthorLastName,
		&cwa.AuthorAvatarURL, &cwa.AuthorEmail,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cwa, nil
}

// ListByIssue retrieves threaded comments for an issue with pagination on top-level comments
func (r *CommentRepository) ListByIssue(ctx context.Context, issueID uuid.UUID, page, limit int) ([]models.CommentThread, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Count total top-level comments
	countQuery := `
		SELECT COUNT(*)
		FROM comments
		WHERE issue_id = $1 AND parent_comment_id IS NULL
	`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, issueID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch top-level comments (paginated)
	topLevelQuery := `
		SELECT c.id, c.issue_id, c.author_id, c.body, c.parent_comment_id,
		       c.is_edited, c.edited_at, c.is_deleted, c.deleted_at, c.created_at, c.updated_at,
		       u.first_name AS author_first_name, u.last_name AS author_last_name,
		       u.avatar_url AS author_avatar_url, u.email AS author_email
		FROM comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.issue_id = $1 AND c.parent_comment_id IS NULL
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, topLevelQuery, issueID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topLevelComments []models.CommentWithAuthor
	var topLevelIDs []uuid.UUID

	for rows.Next() {
		var cwa models.CommentWithAuthor
		if err := rows.Scan(
			&cwa.ID, &cwa.IssueID, &cwa.AuthorID, &cwa.Body,
			&cwa.ParentCommentID, &cwa.IsEdited, &cwa.EditedAt,
			&cwa.IsDeleted, &cwa.DeletedAt, &cwa.CreatedAt, &cwa.UpdatedAt,
			&cwa.AuthorFirstName, &cwa.AuthorLastName,
			&cwa.AuthorAvatarURL, &cwa.AuthorEmail,
		); err != nil {
			return nil, 0, err
		}
		topLevelComments = append(topLevelComments, cwa)
		topLevelIDs = append(topLevelIDs, cwa.ID)
	}

	if len(topLevelComments) == 0 {
		return []models.CommentThread{}, total, nil
	}

	// Fetch all replies for the top-level comments in one query
	repliesQuery := `
		SELECT c.id, c.issue_id, c.author_id, c.body, c.parent_comment_id,
		       c.is_edited, c.edited_at, c.is_deleted, c.deleted_at, c.created_at, c.updated_at,
		       u.first_name AS author_first_name, u.last_name AS author_last_name,
		       u.avatar_url AS author_avatar_url, u.email AS author_email
		FROM comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.parent_comment_id = ANY($1)
		ORDER BY c.created_at ASC
	`

	replyRows, err := r.db.Query(ctx, repliesQuery, topLevelIDs)
	if err != nil {
		return nil, 0, err
	}
	defer replyRows.Close()

	// Group replies by parent_comment_id
	repliesByParent := make(map[uuid.UUID][]models.CommentWithAuthor)
	for replyRows.Next() {
		var reply models.CommentWithAuthor
		if err := replyRows.Scan(
			&reply.ID, &reply.IssueID, &reply.AuthorID, &reply.Body,
			&reply.ParentCommentID, &reply.IsEdited, &reply.EditedAt,
			&reply.IsDeleted, &reply.DeletedAt, &reply.CreatedAt, &reply.UpdatedAt,
			&reply.AuthorFirstName, &reply.AuthorLastName,
			&reply.AuthorAvatarURL, &reply.AuthorEmail,
		); err != nil {
			return nil, 0, err
		}
		if reply.ParentCommentID != nil {
			repliesByParent[*reply.ParentCommentID] = append(repliesByParent[*reply.ParentCommentID], reply)
		}
	}

	// Build threaded response
	threads := make([]models.CommentThread, len(topLevelComments))
	for i, comment := range topLevelComments {
		replies := repliesByParent[comment.ID]
		if replies == nil {
			replies = []models.CommentWithAuthor{}
		}
		threads[i] = models.CommentThread{
			Comment: comment,
			Replies: replies,
		}
	}

	return threads, total, nil
}

// Update updates a comment's body and marks it as edited
func (r *CommentRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateCommentParams) (*models.Comment, error) {
	query := `
		UPDATE comments
		SET body = $2, is_edited = true, edited_at = NOW()
		WHERE id = $1 AND is_deleted = false
		RETURNING id, issue_id, author_id, body, parent_comment_id,
		          is_edited, edited_at, is_deleted, deleted_at, created_at, updated_at
	`

	var comment models.Comment
	err := r.db.QueryRow(ctx, query, id, params.Body).Scan(
		&comment.ID, &comment.IssueID, &comment.AuthorID, &comment.Body,
		&comment.ParentCommentID, &comment.IsEdited, &comment.EditedAt,
		&comment.IsDeleted, &comment.DeletedAt, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &comment, nil
}

// SoftDelete marks a comment as deleted and clears its body
func (r *CommentRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE comments
		SET is_deleted = true, deleted_at = NOW(), body = '[deleted]'
		WHERE id = $1 AND is_deleted = false
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// CountByIssue returns the count of non-deleted comments for an issue
func (r *CommentRepository) CountByIssue(ctx context.Context, issueID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE issue_id = $1 AND is_deleted = false`
	var count int
	err := r.db.QueryRow(ctx, query, issueID).Scan(&count)
	return count, err
}

// GetAuthorID returns the author_id for a comment (used for authorization)
func (r *CommentRepository) GetAuthorID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	query := `SELECT author_id FROM comments WHERE id = $1`
	var authorID uuid.UUID
	err := r.db.QueryRow(ctx, query, id).Scan(&authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &authorID, nil
}
