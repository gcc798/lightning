package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	"github.com/gcc798/lightning/internal/config"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/registry"
	"github.com/gcc798/lightning/internal/transport"
	"github.com/google/uuid"
)

// @title Lightning API
// @version 1.0
// @description Lightning RESTful API
// @host localhost:9009
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg, _, err := config.Load("application/gateway", config.ServiceGateway)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	log, err := logging.NewLogger(config.CurrentEnv(), cfg.AppDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	reg, err := registry.New(cfg.Registry.Driver, cfg.Registry.Address, cfg.Registry.Prefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	defer reg.Close()
	pool := transport.NewClientPool(reg)
	defer pool.Close()
	security := iamv1.NewCached(iamv1.NewRemote(pool), 5*time.Second)
	handler := withGatewayMiddleware(&gateway{registry: reg, selector: registry.NewSelector(), security: security}, cfg)
	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: handler, ReadHeaderTimeout: 10 * time.Second, MaxHeaderBytes: 1 << 20}
	errorsCh := make(chan error, 1)
	go func() {
		var serveErr error
		if cfg.Server.TLSCertFile != "" || cfg.Server.TLSKeyFile != "" {
			if cfg.Server.TLSCertFile == "" || cfg.Server.TLSKeyFile == "" {
				serveErr = errors.New("both server.tlsCertFile and server.tlsKeyFile are required")
			} else {
				serveErr = srv.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
			}
		} else {
			serveErr = srv.ListenAndServe()
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errorsCh <- serveErr
		}
	}()
	log.Info("gateway started")
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
			return
		}
	case err := <-errorsCh:
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}

type rateLimiter struct {
	mu     sync.Mutex
	window time.Time
	counts map[string]int
	limit  int
}

func withGatewayMiddleware(next http.Handler, cfg *config.Config) http.Handler {
	limiter := &rateLimiter{window: time.Now(), counts: make(map[string]int), limit: cfg.Gateway.RateLimitPerMinute}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
			request.Header.Set("X-Request-ID", requestID)
		}
		writer.Header().Set("X-Request-ID", requestID)
		if cfg.CORS.Enabled {
			writer.Header().Set("Access-Control-Allow-Origin", "*")
			writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, clientid, X-Request-ID")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
		}
		host, _, _ := net.SplitHostPort(request.RemoteAddr)
		if !limiter.allow(host) {
			writeError(writer, http.StatusTooManyRequests, "请求过于频繁")
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (l *rateLimiter) allow(key string) bool {
	if l.limit <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.window) >= time.Minute {
		l.window = time.Now()
		clear(l.counts)
	}
	l.counts[key]++
	return l.counts[key] <= l.limit
}
