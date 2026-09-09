package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func TestRouterPreservesMiddlewareBeforeController(t *testing.T) {
	router := NewRouter(echo.New())
	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("authorized", true)
			return next(c)
		}
	}
	router.GET("/test", middleware, func(c *echo.Context) {
		require.Equal(t, true, c.Get("authorized"))
		_ = c.NoContent(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/test", nil))
	require.Equal(t, http.StatusNoContent, recorder.Code)
}
