package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment on an issue
type Comment struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	IssueID         uuid.UUID  `json:"issue_id" db:"issue_id"`
	AuthorID        uuid.UUID  `json:"author_id" db:"author_id"`
	Body            string     `json:"body" db:"body"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" db:"parent_comment_id"`
	IsEdited        bool       `json:"is_edited" db:"is_edited"`
	EditedAt        *time.Time `json:"edited_at" db:"edited_at"`
	IsDeleted       bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// CommentWithAuthor extends Comment with author details from the users table
type CommentWithAuthor struct {
	Comment
	AuthorFirstName *string `json:"author_first_name" db:"author_first_name"`
	AuthorLastName  *string `json:"author_last_name" db:"author_last_name"`
	AuthorAvatarURL *string `json:"author_avatar_url" db:"author_avatar_url"`
	AuthorEmail     string  `json:"author_email" db:"author_email"`
}

// CommentThread represents a top-level comment with its replies
type CommentThread struct {
	Comment CommentWithAuthor   `json:"comment"`
	Replies []CommentWithAuthor `json:"replies"`
}

// CreateCommentParams contains the fields needed to create a comment
type CreateCommentParams struct {
	IssueID         uuid.UUID  `json:"issue_id"`
	AuthorID        uuid.UUID  `json:"author_id"`
	Body            string     `json:"body"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
}

// UpdateCommentParams contains the fields that can be updated on a comment
type UpdateCommentParams struct {
	Body string `json:"body"`
}
