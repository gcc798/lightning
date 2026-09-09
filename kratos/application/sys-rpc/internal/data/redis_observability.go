package data

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	appmetrics "github.com/gcc798/lightning/kratos/pkg/metrics"
	"github.com/redis/go-redis/v9"
)

type redisObservability struct {
	slowThreshold time.Duration
}

type redisMetricsHook struct {
	obs redisObservability
}

var _ redis.Hook = (*redisMetricsHook)(nil)

func newRedisMetricsHook(obs redisObservability) *redisMetricsHook {
	return &redisMetricsHook{obs: obs}
}

func (h *redisMetricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h *redisMetricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		h.obs.observeRedis(cmd.Name(), cmd.Args(), start, err)
		return err
	}
}

func (h *redisMetricsHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		h.obs.observeRedis("pipeline", pipelineArgs(cmds), start, err)
		return err
	}
}

func (o redisObservability) observeRedis(operation string, args []any, start time.Time, err error) {
	duration := time.Since(start)
	operation = strings.ToLower(strings.TrimSpace(operation))
	if operation == "" {
		operation = "unknown"
	}
	status := observeStatus(err)
	appmetrics.ObserveRedis(operation, status, duration)
	if o.slowThreshold > 0 && duration >= o.slowThreshold {
		log.Printf("level=WARN msg=%q operation=%s duration_ms=%d status=%s redis=%v err=%v",
			"slow redis operation",
			operation,
			duration.Milliseconds(),
			status,
			args,
			err,
		)
	}
}

func pipelineArgs(cmds []redis.Cmder) []any {
	items := make([]any, 0, len(cmds))
	for _, cmd := range cmds {
		items = append(items, cmd.Args())
	}
	return items
}
