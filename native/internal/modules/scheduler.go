package modules

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	pkgscheduler "github.com/gcc798/lightning/internal/platform/scheduler"
	"github.com/gcc798/lightning/internal/runtimeconfig"
	"go.uber.org/zap"
)

type SchedulerModule struct {
	mu      sync.Mutex
	cont    Container
	source  runtimeconfig.Source
	jobs    map[string]func(context.Context)
	digest  [sha256.Size]byte
	config  runtimeconfig.SchedulerConfig
	active  *pkgscheduler.Scheduler
	cancel  context.CancelFunc
	done    chan struct{}
	started bool
}

func NewSchedulerModule(jobs map[string]func(context.Context)) *SchedulerModule {
	copyJobs := make(map[string]func(context.Context), len(jobs))
	for name, job := range jobs {
		copyJobs[name] = job
	}
	return &SchedulerModule{jobs: copyJobs}
}

func (*SchedulerModule) Name() string { return SchedulerName }

func (m *SchedulerModule) Init(ctx context.Context, cont Container) error {
	m.cont = cont
	if m.source == nil {
		m.source = runtimeconfig.NewGormSource(cont.GetDB())
	}
	_, cfg, digest, err := m.load(ctx)
	if err != nil {
		return err
	}
	if err := m.validateJobs(cfg); err != nil {
		return err
	}
	m.config, m.digest = cfg, digest
	return nil
}

func (m *SchedulerModule) Start(context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return errors.New("scheduler module is already started")
	}
	if err := m.applyLocked(m.config); err != nil {
		return err
	}
	loopCtx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.done = make(chan struct{})
	m.started = true
	go m.poll(loopCtx, m.done)
	return nil
}

func (m *SchedulerModule) Stop(context.Context) error {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return nil
	}
	cancel, done := m.cancel, m.done
	m.started = false
	m.cancel, m.done = nil, nil
	m.mu.Unlock()

	cancel()
	<-done

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active == nil {
		return nil
	}
	err := m.active.Stop()
	m.active = nil
	return err
}

func (m *SchedulerModule) Refresh(ctx context.Context, req ModuleRefreshRequest) error {
	if !containsCode(req.Codes, runtimeconfig.CodeScheduler) {
		return nil
	}
	_, cfg, digest, err := m.load(ctx)
	if err != nil {
		return err
	}
	if err := m.validateJobs(cfg); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if digest == m.digest {
		return nil
	}
	if m.started {
		if err := m.applyLocked(cfg); err != nil {
			return err
		}
	}
	m.config, m.digest = cfg, digest
	return nil
}

func (m *SchedulerModule) poll(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	for {
		m.mu.Lock()
		seconds := m.config.RefreshIntervalSeconds
		m.mu.Unlock()
		if seconds <= 0 {
			seconds = 5
		}
		timer := time.NewTimer(time.Duration(seconds) * time.Second)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			if err := m.Refresh(ctx, ModuleRefreshRequest{Codes: []string{runtimeconfig.CodeScheduler}, Reason: "poll"}); err != nil {
				m.cont.GetLogger().Error("refresh scheduler configuration failed", zap.Error(err))
			}
		}
	}
}

func (m *SchedulerModule) load(ctx context.Context) ([]byte, runtimeconfig.SchedulerConfig, [sha256.Size]byte, error) {
	raw, err := m.source.Load(ctx, runtimeconfig.CodeScheduler)
	if err != nil {
		return nil, runtimeconfig.SchedulerConfig{}, [sha256.Size]byte{}, err
	}
	var cfg runtimeconfig.SchedulerConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, cfg, [sha256.Size]byte{}, fmt.Errorf("decode scheduler configuration: %w", err)
	}
	return raw, cfg, sha256.Sum256(raw), nil
}

func (m *SchedulerModule) validateJobs(cfg runtimeconfig.SchedulerConfig) error {
	for name, job := range cfg.Jobs {
		if _, exists := m.jobs[name]; job.Enabled && !exists {
			return fmt.Errorf("scheduler job %q is not registered in code", name)
		}
	}
	return nil
}

func (m *SchedulerModule) applyLocked(cfg runtimeconfig.SchedulerConfig) error {
	var candidate *pkgscheduler.Scheduler
	if cfg.Enabled {
		candidate = pkgscheduler.New(m.cont.GetLogger())
		for name, jobConfig := range cfg.Jobs {
			if !jobConfig.Enabled {
				continue
			}
			if err := candidate.AddJob(jobConfig.Cron, name, m.jobs[name]); err != nil {
				return err
			}
		}
	}
	if m.active != nil {
		if err := m.active.Stop(); err != nil {
			return err
		}
	}
	if candidate != nil {
		if err := candidate.Start(); err != nil {
			return err
		}
	}
	m.active = candidate
	return nil
}
