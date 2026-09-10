package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	"github.com/gcc798/microservice-kit/internal/container"
	middleware "github.com/gcc798/microservice-kit/internal/httpmiddleware"
	"github.com/gcc798/microservice-kit/internal/httpx"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/registry"
	"github.com/gcc798/microservice-kit/internal/telemetry"
	"github.com/gcc798/microservice-kit/internal/validator"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	http       *http.Server
	log        logging.Logger
	stopWriter func()
}

// New builds the HTTP server and returns its public business routes for service registration.
func New(cont container.Container, systemAPI sysv1.API, setup func(*httpx.Router) error) (*Server, []registry.HTTPRoute, error) {
	e := echo.New()
	validator.Init()
	e.Validator = validator.EchoValidator{}
	e.Binder = &validator.EchoBinder{}
	e.Use(
		echootel.NewMiddlewareWithConfig(echootel.Config{
			ServerName: fmt.Sprintf("0.0.0.0:%d", cont.GetConfig().Server.Port),
			Skipper: func(c *echo.Context) bool {
				return !telemetry.TraceHTTPPath(c.Request().URL.Path)
			},
		}),
		middleware.Recovery(cont.GetLogger()),
		echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
			LogLatency: true, LogRemoteIP: true, LogMethod: true, LogURIPath: true,
			LogRoutePath: true, LogRequestID: true, LogStatus: true,
			LogValuesFunc: func(c *echo.Context, values echoMiddleware.RequestLoggerValues) error {
				fields := []zap.Field{
					zap.String("method", values.Method),
					zap.String("path", values.URIPath),
					zap.String("route", values.RoutePath),
					zap.Int("status", values.Status),
					zap.Duration("latency", values.Latency),
					zap.String("remote_ip", values.RemoteIP),
					zap.String("request_id", values.RequestID),
				}
				if values.Error != nil {
					fields = append(fields, zap.Error(values.Error))
				}
				logging.WithContext(c.Request().Context(), cont.GetLogger()).Info("http request", fields...)
				return nil
			},
		}),
	)
	if cont.GetConfig().CORS.Enabled {
		e.Use(middleware.CORS())
	}
	e.Use(middleware.StringIDConverter())
	stopWriter := func() {}
	if systemAPI != nil {
		writer := middleware.NewOperLogWriter(systemAPI, cont.GetLogger())
		stopWriter = writer.Stop
		e.Use(middleware.OperationLog(writer))
	}
	if err := setup(httpx.NewRouter(e)); err != nil {
		stopWriter()
		return nil, nil, err
	}
	addr := fmt.Sprintf(":%d", cont.GetConfig().Server.Port)
	srv := &http.Server{Addr: addr, Handler: e, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, MaxHeaderBytes: 1 << 20}
	routes := make([]registry.HTTPRoute, 0, len(e.Router().Routes()))
	for _, route := range e.Router().Routes() {
		if route.Path == "/metrics" || route.Path == "/health" || strings.HasPrefix(route.Path, "/health/") {
			continue
		}
		routes = append(routes, registry.HTTPRoute{Method: route.Method, Path: route.Path})
	}
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})
	return &Server{http: srv, log: cont.GetLogger(), stopWriter: stopWriter}, routes, nil
}

func (s *Server) Run(ctx context.Context) error {
	defer s.stopWriter()
	errs := make(chan error, 1)
	go func() {
		s.log.Info("starting http server", zap.String("addr", s.http.Addr))
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return s.http.Shutdown(shutdown)
	case err := <-errs:
		return err
	}
}
