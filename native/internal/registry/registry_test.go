package registry

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestInProcessResolveWatchAndSelect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reg := NewInProcess()
	defer reg.Close()
	watcher, err := reg.Watch(ctx, "iam")
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Stop()
	if err := reg.Register(ctx, ServiceInstance{ID: "iam-1", Name: "iam", Endpoints: map[string]string{EndpointGRPC: "127.0.0.1:1"}}); err != nil {
		t.Fatal(err)
	}
	instances, err := watcher.Next(ctx)
	if err != nil || len(instances) != 1 {
		t.Fatalf("instances=%v err=%v", instances, err)
	}
	selected, err := NewSelector().Pick("iam", instances)
	if err != nil || selected.ID != "iam-1" {
		t.Fatalf("selected=%v err=%v", selected, err)
	}
}

func TestInProcessWatcherCanStopAfterRegistryClose(t *testing.T) {
	reg := NewInProcess()
	watcher, err := reg.Watch(context.Background(), "iam")
	if err != nil {
		t.Fatal(err)
	}
	if err := reg.Close(); err != nil {
		t.Fatal(err)
	}
	watcher.Stop()
}

func TestConsulRegistryIntegration(t *testing.T) {
	address := os.Getenv("MS_K_TEST_CONSUL_ADDRESS")
	if address == "" {
		t.Skip("MS_K_TEST_CONSUL_ADDRESS is not set")
	}
	reg := NewConsul(address)
	defer reg.Close()
	testRegistryLifecycle(t, reg)
}

func TestEtcdRegistryIntegration(t *testing.T) {
	address := os.Getenv("MS_K_TEST_ETCD_ADDRESS")
	if address == "" {
		t.Skip("MS_K_TEST_ETCD_ADDRESS is not set")
	}
	reg := NewEtcd(address, "/microservice-kit/integration-test")
	defer reg.Close()
	testRegistryLifecycle(t, reg)

	instance := ServiceInstance{ID: "keepalive", Name: "iam", Endpoints: map[string]string{EndpointGRPC: "127.0.0.1:9110"}}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := reg.Register(ctx, instance); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = reg.Deregister(context.Background(), instance) })
	time.Sleep(16 * time.Second)
	instances, err := reg.Resolve(ctx, instance.Name)
	if err != nil || !containsInstance(instances, instance.ID) {
		t.Fatalf("lease was not kept alive: instances=%v err=%v", instances, err)
	}
}

func TestNacosRegistryIntegration(t *testing.T) {
	address := os.Getenv("MS_K_TEST_NACOS_ADDRESS")
	if address == "" {
		t.Skip("MS_K_TEST_NACOS_ADDRESS is not set")
	}
	reg, err := NewNacos(Options{Address: address})
	if err != nil {
		t.Fatal(err)
	}
	defer reg.Close()
	testRegistryLifecycle(t, reg)
}

func TestInstanceMetadataRoundTrip(t *testing.T) {
	want := ServiceInstance{
		ID: "iam-1", Endpoints: map[string]string{EndpointHTTP: "http://127.0.0.1:9010", EndpointGRPC: "127.0.0.1:9110"},
		Routes: []HTTPRoute{{Method: "GET", Path: "/api/v1/user/:id"}, {Method: "POST", Path: "/api/v1/user/page"}},
	}
	metadata, err := encodeMetadata(want)
	if err != nil {
		t.Fatal(err)
	}
	id, endpoints, routes, err := decodeMetadata(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if id != want.ID || endpoints[EndpointHTTP] != want.Endpoints[EndpointHTTP] || len(routes) != len(want.Routes) || routes[0] != want.Routes[0] {
		t.Fatalf("metadata round trip = id %q endpoints %v routes %v", id, endpoints, routes)
	}
}

func TestParseNacosServers(t *testing.T) {
	servers, err := parseNacosServers("http://nacos-1:8848/nacos,https://nacos-2:9443/custom")
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 || servers[0].IpAddr != "nacos-1" || servers[0].ContextPath != "/nacos" || servers[1].Scheme != "https" || servers[1].Port != 9443 {
		t.Fatalf("servers = %+v", servers)
	}
}

func testRegistryLifecycle(t *testing.T, reg Registry) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	instance := ServiceInstance{ID: "lifecycle", Name: "iam", Endpoints: map[string]string{EndpointGRPC: "127.0.0.1:9110"}}
	if err := reg.Register(ctx, instance); err != nil {
		t.Fatal(err)
	}
	instances, err := reg.Resolve(ctx, instance.Name)
	if err != nil || !containsInstance(instances, instance.ID) {
		t.Fatalf("registered instance not resolved: instances=%v err=%v", instances, err)
	}
	if err := reg.Deregister(ctx, instance); err != nil {
		t.Fatal(err)
	}
	instances, err = reg.Resolve(ctx, instance.Name)
	if err != nil || containsInstance(instances, instance.ID) {
		t.Fatalf("deregistered instance still resolved: instances=%v err=%v", instances, err)
	}
}

func containsInstance(instances []ServiceInstance, id string) bool {
	for _, instance := range instances {
		if instance.ID == id {
			return true
		}
	}
	return false
}
