package registry

import (
	"context"
	"sync"
)

// InProcess is a unit-test-only registry; its state cannot be shared across processes.
type InProcess struct {
	mu       sync.RWMutex
	services map[string]map[string]ServiceInstance
	watchers map[string]map[chan struct{}]struct{}
	closed   bool
}

func NewInProcess() *InProcess {
	return &InProcess{services: make(map[string]map[string]ServiceInstance), watchers: make(map[string]map[chan struct{}]struct{})}
}

func (r *InProcess) Register(_ context.Context, instance ServiceInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return ErrWatcherStopped
	}
	if r.services[instance.Name] == nil {
		r.services[instance.Name] = make(map[string]ServiceInstance)
	}
	r.services[instance.Name][instance.ID] = instance
	r.notify(instance.Name)
	return nil
}

func (r *InProcess) Deregister(_ context.Context, instance ServiceInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services[instance.Name], instance.ID)
	r.notify(instance.Name)
	return nil
}

func (r *InProcess) Resolve(_ context.Context, name string) ([]ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instances := make([]ServiceInstance, 0, len(r.services[name]))
	for _, instance := range r.services[name] {
		instances = append(instances, instance)
	}
	return instances, nil
}

func (r *InProcess) Watch(_ context.Context, name string) (Watcher, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	updates := make(chan struct{}, 1)
	if r.watchers[name] == nil {
		r.watchers[name] = make(map[chan struct{}]struct{})
	}
	r.watchers[name][updates] = struct{}{}
	updates <- struct{}{}
	return &inProcessWatcher{registry: r, name: name, updates: updates}, nil
}

func (r *InProcess) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	r.closed = true
	for _, watchers := range r.watchers {
		for updates := range watchers {
			close(updates)
		}
	}
	clear(r.watchers)
	return nil
}

func (r *InProcess) notify(name string) {
	for updates := range r.watchers[name] {
		select {
		case updates <- struct{}{}:
		default:
		}
	}
}

type inProcessWatcher struct {
	registry *InProcess
	name     string
	updates  chan struct{}
	once     sync.Once
}

func (w *inProcessWatcher) Next(ctx context.Context) ([]ServiceInstance, error) {
	select {
	case _, ok := <-w.updates:
		if !ok {
			return nil, ErrWatcherStopped
		}
		return w.registry.Resolve(ctx, w.name)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (w *inProcessWatcher) Stop() {
	w.once.Do(func() {
		w.registry.mu.Lock()
		defer w.registry.mu.Unlock()
		if _, ok := w.registry.watchers[w.name][w.updates]; ok {
			delete(w.registry.watchers[w.name], w.updates)
			close(w.updates)
		}
	})
}
