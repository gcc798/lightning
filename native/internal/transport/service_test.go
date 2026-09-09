package transport

import (
	"context"
	"strings"
	"testing"

	"github.com/gcc798/lightning/internal/registry"
)

func TestRegisterService(t *testing.T) {
	reg := registry.NewInProcess()
	defer reg.Close()

	instance, err := RegisterService(context.Background(), reg, "scheduler", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(instance.ID, "scheduler-") {
		t.Fatalf("instance ID = %q", instance.ID)
	}
	instances, err := reg.Resolve(context.Background(), "scheduler")
	if err != nil || len(instances) != 1 || instances[0].ID != instance.ID {
		t.Fatalf("instances = %v, err = %v", instances, err)
	}
}
