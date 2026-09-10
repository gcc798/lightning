package websocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func TestOriginAllowed(t *testing.T) {
	for _, test := range []struct {
		name, origin, host string
		allow, want        bool
	}{
		{name: "no browser origin", host: "api.example.com", want: true},
		{name: "same origin", origin: "https://api.example.com", host: "api.example.com", want: true},
		{name: "cross origin blocked", origin: "https://evil.example.com", host: "api.example.com"},
		{name: "development override", origin: "http://localhost:5173", host: "localhost:9009", allow: true, want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := &http.Request{Host: test.host, Header: http.Header{"Origin": []string{test.origin}}}
			if got := originAllowed(req, test.allow); got != test.want {
				t.Fatalf("originAllowed() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHandlerUpgradesEchoResponse(t *testing.T) {
	log, err := logging.NewLoggerWithConfig(&logging.Config{Level: "error", Output: "console", Encoding: "json"})
	require.NoError(t, err)
	hub := NewHub(log)
	hub.Start()
	t.Cleanup(hub.Close)

	e := echo.New()
	handler := NewHandler(hub, log, false)
	e.GET("/resource/websocket", func(c *echo.Context) error {
		handler.ServeWs(c)
		return nil
	})
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/resource/websocket?userId=42"
	conn, response, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	require.Equal(t, 101, response.StatusCode)
	require.NoError(t, conn.Close())
}
