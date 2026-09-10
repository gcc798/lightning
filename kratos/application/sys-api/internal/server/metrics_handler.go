package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	appmetrics "github.com/gcc798/microservice-kit/kratos/pkg/metrics"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerMetricsEndpoint(srv *khttp.Server) {
	if srv == nil {
		return
	}
	handler := promhttp.Handler()
	srv.Route("/").GET("/metrics", wrapOperation(operationMetrics, func(ctx context.Context, httpCtx khttp.Context) error {
		handler.ServeHTTP(httpCtx.Response(), httpCtx.Request())
		return nil
	}))
}

func metricsMiddleware() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			method := ""
			path := ""
			if httpReq, ok := khttp.RequestFromServerContext(ctx); ok && httpReq != nil {
				method = httpReq.Method
				if httpReq.URL != nil {
					path = httpReq.URL.Path
				}
			}
			startedAt := time.Now()
			reply, err := next(ctx, req)
			status := http.StatusOK
			if method == "" {
				method = http.MethodGet
			}
			if path == "" {
				path = "/"
			}
			appmetrics.HTTPRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
			appmetrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(time.Since(startedAt).Seconds())
			return reply, err
		}
	}
}
