package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestResponseHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
		call func(*echo.Context)
	}{
		{name: "success", want: http.StatusOK, call: func(c *echo.Context) { Success(c, nil) }},
		{name: "bad request", want: http.StatusBadRequest, call: func(c *echo.Context) { BadRequest(c, "bad") }},
		{name: "unauthorized", want: http.StatusUnauthorized, call: func(c *echo.Context) { Unauthorized(c, "no") }},
		{name: "forbidden", want: http.StatusForbidden, call: func(c *echo.Context) { Forbidden(c, "no") }},
		{name: "not found", want: http.StatusNotFound, call: func(c *echo.Context) { NotFound(c, "missing") }},
		{name: "server error", want: http.StatusInternalServerError, call: func(c *echo.Context) { Fail(c, "failed") }},
		{name: "business error", want: http.StatusOK, call: func(c *echo.Context) { FailCode(c, 10001, "business") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			context := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), recorder)
			tt.call(context)
			if recorder.Code != tt.want {
				t.Fatalf("HTTP status = %d, want %d", recorder.Code, tt.want)
			}
		})
	}
}
