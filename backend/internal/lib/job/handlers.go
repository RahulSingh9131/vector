package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RahulSingh9131/vector/internal/config"
	"github.com/RahulSingh9131/vector/internal/database"
	"github.com/RahulSingh9131/vector/internal/lib/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger, db database.DBTX) {
	emailClient = email.NewClient(config, logger)
	j.db = db
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Processing welcome email task")

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")
	return nil
}

func (j *JobService) handleActivityCleanupTask(ctx context.Context, t *asynq.Task) error {
	cutoff := time.Now().AddDate(0, 0, -ActivityRetentionDays)
	totalDeleted := int64(0)

	j.logger.Info().
		Time("cutoff", cutoff).
		Int("retention_days", ActivityRetentionDays).
		Msg("Starting activity cleanup")

	for {
		query := `
			DELETE FROM activities
			WHERE id IN (
				SELECT id FROM activities
				WHERE created_at < $1
				LIMIT $2
			)
		`

		result, err := j.db.Exec(ctx, query, cutoff, CleanupBatchSize)
		if err != nil {
			j.logger.Error().Err(err).
				Int64("deleted_so_far", totalDeleted).
				Msg("Activity cleanup batch failed")
			return fmt.Errorf("cleanup batch failed after deleting %d rows: %w", totalDeleted, err)
		}

		deleted := result.RowsAffected()
		totalDeleted += deleted

		j.logger.Debug().
			Int64("batch_deleted", deleted).
			Int64("total_deleted", totalDeleted).
			Msg("Activity cleanup batch complete")

		if deleted < int64(CleanupBatchSize) {
			break
		}
	}

	j.logger.Info().
		Int64("total_deleted", totalDeleted).
		Msg("Activity cleanup finished")

	return nil
}
