package modules

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gcc798/microservice-kit/internal/runtimeconfig"
)

func TestSchedulerModuleRunsAndStopsRegisteredJob(t *testing.T) {
	source := &mapSource{data: map[string][]byte{
		runtimeconfig.CodeScheduler: []byte(`{"enabled":true,"refreshIntervalSeconds":1,"jobs":{"tick":{"enabled":true,"cron":"*/1 * * * * *"}}}`),
	}}
	cont := &moduleTestContainer{modules: make(map[string]Module)}
	var runs atomic.Int32
	schedulerModule := NewSchedulerModule(map[string]func(context.Context){"tick": func(context.Context) { runs.Add(1) }})
	schedulerModule.source = source
	cont.modules[SchedulerName] = schedulerModule
	ctx := context.Background()
	if err := schedulerModule.Init(ctx, cont); err != nil {
		t.Fatal(err)
	}
	if err := schedulerModule.Start(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1200 * time.Millisecond)
	if runs.Load() == 0 {
		t.Fatal("registered job did not run")
	}
	disabled := []byte(`{"enabled":false,"refreshIntervalSeconds":1,"jobs":{"tick":{"enabled":true,"cron":"*/1 * * * * *"}}}`)
	source.mu.Lock()
	source.data[runtimeconfig.CodeScheduler] = disabled
	source.mu.Unlock()
	if err := schedulerModule.Refresh(ctx, ModuleRefreshRequest{Codes: []string{runtimeconfig.CodeScheduler}, Reason: "test"}); err != nil {
		t.Fatal(err)
	}
	afterDisable := runs.Load()
	time.Sleep(1100 * time.Millisecond)
	if got := runs.Load(); got != afterDisable {
		t.Fatalf("job continued after runtime disable: before=%d after=%d", afterDisable, got)
	}
	if err := schedulerModule.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	afterStop := runs.Load()
	time.Sleep(1100 * time.Millisecond)
	if got := runs.Load(); got != afterStop {
		t.Fatalf("job continued after stop: before=%d after=%d", afterStop, got)
	}
}
