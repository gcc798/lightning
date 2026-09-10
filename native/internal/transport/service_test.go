package transport

import (
	"context"
	"testing"

	"github.com/gcc798/microservice-kit/internal/registry"
)

func TestRegisterService(t *testing.T) {
	reg := registry.NewInProcess()
	defer reg.Close()

	instance, err := RegisterService(context.Background(), reg, "scheduler", "scheduler-test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if instance.ID != "scheduler-test" {
		t.Fatalf("instance ID = %q", instance.ID)
	}
	instances, err := reg.Resolve(context.Background(), "scheduler")
	if err != nil || len(instances) != 1 || instances[0].ID != instance.ID {
		t.Fatalf("instances = %v, err = %v", instances, err)
	}
	if _, err := RegisterService(context.Background(), reg, "scheduler", "", nil); err == nil {
		t.Fatal("RegisterService accepted an empty instance ID")
	}
}
