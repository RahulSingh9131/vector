package job

const (
	// TaskActivityCleanup is the task type for the periodic activity cleanup job.
	TaskActivityCleanup = "task:activity.cleanup"

	// ActivityRetentionDays is the number of days to retain activity records.
	ActivityRetentionDays = 90

	// CleanupBatchSize is the number of rows to delete per batch to avoid table locks.
	CleanupBatchSize = 1000
)
