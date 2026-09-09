package scheduler

import (
	"context"
	"fmt"
	"sync"

	logging "github.com/gcc798/lightning/internal/logger"

	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Scheduler 调度器管理所有定时任务。
type Scheduler struct {
	cron    *cron.Cron
	logger  logging.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	jobs    map[string]cron.EntryID // 任务名称到 EntryID 的映射
	jobsMux sync.RWMutex            // 保护jobs map的并发访问
}

// New 创建调度器实例。
func New(logger logging.Logger) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:   c,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		jobs:   make(map[string]cron.EntryID),
	}
}

// AddJob 添加命名定时任务。
func (s *Scheduler) AddJob(spec string, name string, job func(context.Context)) error {
	s.jobsMux.Lock()
	defer s.jobsMux.Unlock()

	// 检查任务是否已存在
	if _, exists := s.jobs[name]; exists {
		return fmt.Errorf("job %s already exists", name)
	}

	// 添加任务到cron
	entryID, err := s.cron.AddFunc(spec, func() {
		ctx, span := otel.Tracer("github.com/gcc798/lightning/internal/platform/scheduler").Start(
			s.ctx, "scheduler.job", trace.WithAttributes(attribute.String("job.name", name)),
		)
		defer span.End()
		jobLogger := logging.WithContext(ctx, s.logger)
		jobLogger.Debug("running scheduled job", zap.String("job", name))
		defer func() {
			if r := recover(); r != nil {
				span.RecordError(fmt.Errorf("job panic: %v", r))
				span.SetStatus(codes.Error, "job panic")
				jobLogger.Error("job panic", zap.String("job", name), zap.Any("panic", r))
			}
		}()
		job(ctx)
	})

	if err != nil {
		s.logger.Error("failed to add job", zap.String("job", name), zap.Error(err))
		return err
	}

	// 记录任务ID
	s.jobs[name] = entryID
	s.logger.Info("job added successfully", zap.String("job", name), zap.String("spec", spec))
	return nil
}

// AddFunc 添加匿名定时函数，调用方负责保存返回的 EntryID。
func (s *Scheduler) AddFunc(spec string, job func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, job)
}

// RemoveEntry 按 EntryID 移除定时函数。
func (s *Scheduler) RemoveEntry(entryID cron.EntryID) {
	s.cron.Remove(entryID)
}

// RemoveJob 移除指定名称的定时任务。
func (s *Scheduler) RemoveJob(name string) error {
	s.jobsMux.Lock()
	defer s.jobsMux.Unlock()

	entryID, exists := s.jobs[name]
	if !exists {
		return fmt.Errorf("job %s not found", name)
	}

	// 从cron中移除任务
	s.cron.Remove(entryID)

	// 从映射中删除
	delete(s.jobs, name)

	s.logger.Info("job removed successfully", zap.String("job", name))
	return nil
}

// UpdateJob 更新指定任务的 cron 表达式和执行函数。
func (s *Scheduler) UpdateJob(spec string, name string, job func(context.Context)) error {
	// 先移除旧任务
	if err := s.RemoveJob(name); err != nil {
		// 如果任务不存在，直接添加新任务
		if err.Error() != fmt.Sprintf("job %s not found", name) {
			return err
		}
	}

	// 添加新任务
	return s.AddJob(spec, name, job)
}

// JobExists 检查指定任务是否存在。
func (s *Scheduler) JobExists(name string) bool {
	s.jobsMux.RLock()
	defer s.jobsMux.RUnlock()
	_, exists := s.jobs[name]
	return exists
}

// GetJobCount 返回当前任务数量。
func (s *Scheduler) GetJobCount() int {
	s.jobsMux.RLock()
	defer s.jobsMux.RUnlock()
	return len(s.jobs)
}

// Name 返回组件名称。
func (s *Scheduler) Name() string {
	return "Scheduler"
}

// Start 启动调度器。
func (s *Scheduler) Start() error {
	s.logger.Info("starting scheduler")
	s.cron.Start()
	return nil
}

// Stop 停止调度器。
func (s *Scheduler) Stop() error {
	s.logger.Info("stopping scheduler")
	s.cancel()
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("scheduler stopped")
	return nil
}

// GetContext 获取调度器上下文。
func (s *Scheduler) GetContext() context.Context {
	return s.ctx
}
