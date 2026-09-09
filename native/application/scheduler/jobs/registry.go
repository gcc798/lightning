package jobs

import (
	"context"

	resourcev1 "github.com/gcc798/lightning/internal/api/resource/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	logging "github.com/gcc798/lightning/internal/logger"
	"go.uber.org/zap"
)

// Definitions returns the code-owned jobs that a scheduler process may run.
func Definitions(
	system sysv1.API,
	resource resourcev1.API,
	logger logging.Logger,
) map[string]func(context.Context) {
	definitions := make(map[string]func(context.Context))
	cl := NewDataCleanupJob(system, resource, logger)
	definitions["data-cleanup"] = cl.Run

	logger.Info("scheduler job definitions created", zap.Int("count", len(definitions)))
	return definitions
}
