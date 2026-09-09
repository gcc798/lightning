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
	"github.com/gcc798/lightning/internal/validator"
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
	e.Use(middleware.Recovery(cont.GetLogger().Get()), echoMiddleware.RequestLogger())
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
