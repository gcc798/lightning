package jobs

import (
	"context"

	resourcev1 "github.com/gcc798/microservice-kit/internal/api/resource/v1"
	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"go.uber.org/zap"
)

// DataCleanupJob 定义业务数据结构。
type DataCleanupJob struct {
	system   sysv1.API
	resource resourcev1.API
	logger   logging.Logger
}

// NewDataCleanupJob 创建组件实例。
func NewDataCleanupJob(system sysv1.API, resource resourcev1.API, logger logging.Logger) *DataCleanupJob {
	return &DataCleanupJob{system: system, resource: resource, logger: logger}
}

// Run 执行业务任务。
func (j *DataCleanupJob) Run(ctx context.Context) {
	log := logging.WithContext(ctx, j.logger)
	log.Info("starting data cleanup")

	logs, err := j.system.CleanLogs(ctx, 90)
	if err != nil {
		log.Error("failed to cleanup system logs", zap.Error(err))
	} else {
		log.Info("cleaned up system logs", zap.Int64("loginLogs", logs.LoginLogs), zap.Int64("operationLogs", logs.OperationLogs))
	}
	attachments, err := j.resource.CleanExpired(ctx)
	if err != nil {
		log.Error("failed to cleanup expired attachments", zap.Error(err))
	} else {
		log.Info("cleaned up expired attachments", zap.Int64("cleaned", attachments.Cleaned), zap.Int64("failed", attachments.Failed))
	}

	log.Info("data cleanup completed")
}

// Schedule 返回任务调度表达式。
func (j *DataCleanupJob) Schedule() string { return "0 0 2 * * *" }
