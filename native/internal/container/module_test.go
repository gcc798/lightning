package container

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gcc798/microservice-kit/internal/modules"
)

type fakeModule struct {
	name      string
	mu        *sync.Mutex
	events    *[]string
	initErr   error
	startErr  error
	stopErr   error
	refreshFn func()
}

func (m *fakeModule) Name() string { return m.name }

func (m *fakeModule) record(event string) {
	if m.events == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	*m.events = append(*m.events, event)
}

func (m *fakeModule) Init(context.Context, modules.Container) error {
	m.record("init:" + m.name)
	return m.initErr
}

func (m *fakeModule) Start(context.Context) error {
	m.record("start:" + m.name)
	return m.startErr
}

func (m *fakeModule) Stop(context.Context) error {
	m.record("stop:" + m.name)
	return m.stopErr
}

func (m *fakeModule) Refresh(context.Context, modules.ModuleRefreshRequest) error {
	if m.refreshFn != nil {
		m.refreshFn()
	}
	return nil
}

func TestModuleLifecycleOrder(t *testing.T) {
	var mu sync.Mutex
	events := make([]string, 0)
	c := NewEmpty(nil, nil, nil)
	candidates := []modules.Module{
		&fakeModule{name: "first", mu: &mu, events: &events},
		&fakeModule{name: "second", mu: &mu, events: &events},
	}
	ctx := context.Background()
	if err := c.RegisterModules(ctx, candidates...); err != nil {
		t.Fatal(err)
	}
	if err := c.StartModules(ctx); err != nil {
		t.Fatal(err)
	}
	if err := c.StopModules(ctx); err != nil {
		t.Fatal(err)
	}
	want := []string{"init:first", "init:second", "start:first", "start:second", "stop:second", "stop:first"}
	if !reflect.DeepEqual(events, want) {
		t.Fatalf("events = %#v, want %#v", events, want)
	}
}

func TestStartModulesRollsBackStartedModules(t *testing.T) {
	var mu sync.Mutex
	events := make([]string, 0)
	c := NewEmpty(nil, nil, nil)
	ctx := context.Background()
	if err := c.RegisterModules(ctx,
		&fakeModule{name: "first", mu: &mu, events: &events},
		&fakeModule{name: "second", mu: &mu, events: &events, startErr: errors.New("boom")},
		&fakeModule{name: "third", mu: &mu, events: &events},
	); err != nil {
		t.Fatal(err)
	}
	if err := c.StartModules(ctx); err == nil {
		t.Fatal("StartModules() succeeded")
	}
	wantTail := []string{"start:first", "start:second", "stop:first"}
	if !reflect.DeepEqual(events[len(events)-len(wantTail):], wantTail) {
		t.Fatalf("events tail = %#v, want %#v", events, wantTail)
	}
}

func TestStopModulesAttemptsEveryModule(t *testing.T) {
	var mu sync.Mutex
	events := make([]string, 0)
	c := NewEmpty(nil, nil, nil)
	ctx := context.Background()
	if err := c.RegisterModules(ctx,
		&fakeModule{name: "first", mu: &mu, events: &events, stopErr: errors.New("first stop")},
		&fakeModule{name: "second", mu: &mu, events: &events, stopErr: errors.New("second stop")},
	); err != nil {
		t.Fatal(err)
	}
	if err := c.StartModules(ctx); err != nil {
		t.Fatal(err)
	}
	err := c.StopModules(ctx)
	if err == nil || !containsAll(err.Error(), "first stop", "second stop") {
		t.Fatalf("StopModules() error = %v", err)
	}
	wantTail := []string{"stop:second", "stop:first"}
	if !reflect.DeepEqual(events[len(events)-2:], wantTail) {
		t.Fatalf("events tail = %#v, want %#v", events, wantTail)
	}
}

func TestRefreshModuleIsSerialized(t *testing.T) {
	c := NewEmpty(nil, nil, nil)
	var active atomic.Int32
	var overlap atomic.Bool
	m := &fakeModule{name: "refresh", refreshFn: func() {
		if active.Add(1) != 1 {
			overlap.Store(true)
		}
		time.Sleep(10 * time.Millisecond)
		active.Add(-1)
	}}
	if err := c.RegisterModules(context.Background(), m); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 2)
	for range 2 {
		go func() {
			done <- c.RefreshModule(context.Background(), "refresh", modules.ModuleRefreshRequest{})
		}()
	}
	for range 2 {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
	if overlap.Load() {
		t.Fatal("module refresh calls overlapped")
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}
	return true
}
