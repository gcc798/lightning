package registry

import (
	"context"
	"sync"
	"time"
)

type pollingWatcher struct {
	registry Registry
	name     string
	interval time.Duration
	stop     chan struct{}
	once     sync.Once
	first    bool
}

func newPollingWatcher(registry Registry, name string) Watcher {
	return &pollingWatcher{registry: registry, name: name, interval: 2 * time.Second, stop: make(chan struct{}), first: true}
}

func (w *pollingWatcher) Next(ctx context.Context) ([]ServiceInstance, error) {
	if w.first {
		w.first = false
		return w.registry.Resolve(ctx, w.name)
	}
	timer := time.NewTimer(w.interval)
	defer timer.Stop()
	select {
	case <-timer.C:
		return w.registry.Resolve(ctx, w.name)
	case <-w.stop:
		return nil, ErrWatcherStopped
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (w *pollingWatcher) Stop() { w.once.Do(func() { close(w.stop) }) }
