package repository

import (
	"context"
	"errors"

	models "github.com/RahulSingh9131/vector/internal/model"
	"github.com/RahulSingh9131/vector/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// LabelRepository handles database operations for labels
type LabelRepository struct {
	server *server.Server
}

// NewLabelRepository creates a new label repository
func NewLabelRepository(s *server.Server) *LabelRepository {
	return &LabelRepository{
		server: s,
	}
}

// Create creates a new label
func (r *LabelRepository) Create(ctx context.Context, params models.CreateLabelParams) (*models.Label, error) {
	query := `
		INSERT INTO labels (project_id, name, color, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, project_id, name, color, description, created_at, updated_at
	`

	var label models.Label
	err := r.server.DB.Pool.QueryRow(
		ctx, query,
		params.ProjectID, params.Name, params.Color, params.Description,
	).Scan(
		&label.ID, &label.ProjectID, &label.Name, &label.Color,
		&label.Description, &label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &label, nil
}

// GetByID retrieves a label by its UUID
func (r *LabelRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Label, error) {
	query := `
		SELECT id, project_id, name, color, description, created_at, updated_at
		FROM labels
		WHERE id = $1
	`

	var label models.Label
	err := r.server.DB.Pool.QueryRow(ctx, query, id).Scan(
		&label.ID, &label.ProjectID, &label.Name, &label.Color,
		&label.Description, &label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &label, nil
}

// ListByProject retrieves all labels for a project
func (r *LabelRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.Label, error) {
	query := `
		SELECT id, project_id, name, color, description, created_at, updated_at
		FROM labels
		WHERE project_id = $1
		ORDER BY name ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var label models.Label
		if err := rows.Scan(
			&label.ID, &label.ProjectID, &label.Name, &label.Color,
			&label.Description, &label.CreatedAt, &label.UpdatedAt,
		); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}

// Update updates a label
func (r *LabelRepository) Update(ctx context.Context, id uuid.UUID, params models.UpdateLabelParams) (*models.Label, error) {
	query := `
		UPDATE labels
		SET name = COALESCE($2, name),
		    color = COALESCE($3, color),
		    description = COALESCE($4, description)
		WHERE id = $1
		RETURNING id, project_id, name, color, description, created_at, updated_at
	`

	var label models.Label
	err := r.server.DB.Pool.QueryRow(
		ctx, query,
		id, params.Name, params.Color, params.Description,
	).Scan(
		&label.ID, &label.ProjectID, &label.Name, &label.Color,
		&label.Description, &label.CreatedAt, &label.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &label, nil
}

// Delete deletes a label (cascade removes issue_labels associations)
func (r *LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labels WHERE id = $1`
	_, err := r.server.DB.Pool.Exec(ctx, query, id)
	return err
}

// AddLabelToIssue associates a label with an issue
func (r *LabelRepository) AddLabelToIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	query := `
		INSERT INTO issue_labels (issue_id, label_id)
		VALUES ($1, $2)
		ON CONFLICT (issue_id, label_id) DO NOTHING
	`
	_, err := r.server.DB.Pool.Exec(ctx, query, issueID, labelID)
	return err
}

// RemoveLabelFromIssue removes a label from an issue
func (r *LabelRepository) RemoveLabelFromIssue(ctx context.Context, issueID, labelID uuid.UUID) error {
	query := `DELETE FROM issue_labels WHERE issue_id = $1 AND label_id = $2`
	_, err := r.server.DB.Pool.Exec(ctx, query, issueID, labelID)
	return err
}

// GetLabelsByIssue retrieves all labels attached to an issue
func (r *LabelRepository) GetLabelsByIssue(ctx context.Context, issueID uuid.UUID) ([]models.Label, error) {
	query := `
		SELECT l.id, l.project_id, l.name, l.color, l.description, l.created_at, l.updated_at
		FROM labels l
		INNER JOIN issue_labels il ON il.label_id = l.id
		WHERE il.issue_id = $1
		ORDER BY l.name ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []models.Label
	for rows.Next() {
		var label models.Label
		if err := rows.Scan(
			&label.ID, &label.ProjectID, &label.Name, &label.Color,
			&label.Description, &label.CreatedAt, &label.UpdatedAt,
		); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}
