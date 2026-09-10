package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	iamv1 "github.com/gcc798/microservice-kit/internal/api/iam/v1"
	"github.com/gcc798/microservice-kit/internal/registry"
)

func TestSMSCodeRouteIsPublicIAM(t *testing.T) {
	const path = "/resource/sms/code"
	table, err := buildRouteTable(map[string][]registry.ServiceInstance{
		iamv1.ServiceName: {{ID: "iam-1", Name: iamv1.ServiceName, Routes: []registry.HTTPRoute{{Method: http.MethodGet, Path: path}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	proxy := &gateway{}
	proxy.routes.Store(table)
	target, err := proxy.route(http.MethodGet, path)
	if err != nil || target.service != iamv1.ServiceName {
		t.Fatalf("route(%q) = %+v, %v", path, target, err)
	}
	if !publicPath(path) {
		t.Fatalf("publicPath(%q) = false", path)
	}
}

func TestGatewayStartupHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/startup", nil)
	response := httptest.NewRecorder()
	(&gateway{}).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestGatewayOwnsSwagger(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	response := httptest.NewRecorder()
	(&gateway{}).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if _, err := (&gateway{}).route(request.Method, request.URL.Path); err == nil {
		t.Fatal("swagger was routed to a service")
	}
}

func TestRouteTableMatchesMethodAndRESTPath(t *testing.T) {
	instance := registry.ServiceInstance{
		ID: "resource-1", Name: "resource",
		Routes: []registry.HTTPRoute{{Method: http.MethodGet, Path: "/api/v1/attachment/:id"}},
	}
	table, err := buildRouteTable(map[string][]registry.ServiceInstance{"resource": {instance}})
	if err != nil {
		t.Fatal(err)
	}
	proxy := &gateway{}
	proxy.routes.Store(table)
	if target, err := proxy.route(http.MethodGet, "/api/v1/attachment/42"); err != nil || target.instances[0].ID != instance.ID {
		t.Fatalf("REST route target=%+v err=%v", target, err)
	}
	if _, err := proxy.route(http.MethodPost, "/api/v1/attachment/42"); err == nil {
		t.Fatal("route ignored HTTP method")
	}
}

func TestRouteTablePrefersStaticSegment(t *testing.T) {
	table, err := buildRouteTable(map[string][]registry.ServiceInstance{
		"first":  {{ID: "first-1", Routes: []registry.HTTPRoute{{Method: http.MethodGet, Path: "/items/:id/detail"}}}},
		"second": {{ID: "second-1", Routes: []registry.HTTPRoute{{Method: http.MethodGet, Path: "/items/current/:id"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	proxy := &gateway{}
	proxy.routes.Store(table)
	target, err := proxy.route(http.MethodGet, "/items/current/detail")
	if err != nil || target.service != "second" {
		t.Fatalf("target=%+v err=%v", target, err)
	}
}
