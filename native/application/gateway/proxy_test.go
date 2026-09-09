package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
)

func TestSMSCodeRouteIsPublicIAM(t *testing.T) {
	const path = "/resource/sms/code"
	if service := routeService(path); service != iamv1.ServiceName {
		t.Fatalf("routeService(%q) = %q, want %q", path, service, iamv1.ServiceName)
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
	if service := routeService(request.URL.Path); service != "" {
		t.Fatalf("swagger routed to %q", service)
	}
}
