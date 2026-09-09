package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestCORSDevelopmentRequest(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.Use(CORS())
	e.GET("/resource", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/resource", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	recorder := httptest.NewRecorder()
	e.ServeHTTP(recorder, request)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func TestCORSPreflightStopsHandler(t *testing.T) {
	t.Parallel()

	e := echo.New()
	e.Use(CORS())
	called := false
	e.OPTIONS("/resource", func(c *echo.Context) error {
		called = true
		return c.NoContent(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	recorder := httptest.NewRecorder()
	e.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if called {
		t.Fatal("preflight request reached handler")
	}
}
