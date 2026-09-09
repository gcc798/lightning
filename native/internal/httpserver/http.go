package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/container"
	middleware "github.com/gcc798/lightning/internal/httpmiddleware"
	"github.com/gcc798/lightning/internal/httpx"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/telemetry"
	"github.com/gcc798/lightning/internal/validator"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

// RunHTTP owns the process HTTP lifecycle; each service supplies its own routes.
func RunHTTP(ctx context.Context, cont container.Container, security iamv1.API, systemAPI sysv1.API, setup func(*httpx.Router) error) error {
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
	if systemAPI != nil {
		writer := middleware.NewOperLogWriter(systemAPI, cont.GetLogger())
		defer writer.Stop()
		e.Use(middleware.OperationLog(writer))
	}
	if err := setup(httpx.NewRouter(e)); err != nil {
		return err
	}
	addr := fmt.Sprintf(":%d", cont.GetConfig().Server.Port)
	srv := &http.Server{Addr: addr, Handler: e, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, MaxHeaderBytes: 1 << 20}
	errs := make(chan error, 1)
	go func() {
		cont.GetLogger().Info("starting http server", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errs <- err
		}
	}()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	case err := <-errs:
		return err
	}
}
