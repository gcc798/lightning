package grpcx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gcc798/microservice-kit/kratos/pkg/registryx"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

type ClientConfig struct {
	Mode      string
	Endpoint  string
	Timeout   time.Duration
	Discovery registryx.DiscoveryConfig
}

func DialInsecure(ctx context.Context, cfg ClientConfig, unaryInts ...grpc.UnaryClientInterceptor) (*grpc.ClientConn, func(), error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = "direct"
	}
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	options := []kgrpc.ClientOption{
		kgrpc.WithTimeout(timeout),
		kgrpc.WithOptions(grpc.WithBlock()),
	}
	if len(unaryInts) > 0 {
		options = append(options, kgrpc.WithUnaryInterceptor(unaryInts...))
	}
	switch mode {
	case "direct":
		endpoint := strings.TrimSpace(cfg.Endpoint)
		if endpoint == "" {
			return nil, nil, fmt.Errorf("grpc endpoint is empty in direct mode")
		}
		options = append(options, kgrpc.WithEndpoint(endpoint))
		conn, err := kgrpc.DialInsecure(ctx, options...)
		return conn, nil, err
	case "discovery", "registry":
		discovery, cleanup, err := registryx.NewDiscovery(cfg.Discovery)
		if err != nil {
			return nil, nil, err
		}
		service := strings.TrimSpace(cfg.Discovery.Service)
		if service == "" {
			cleanup()
			return nil, nil, fmt.Errorf("grpc discovery service is empty")
		}
		options = append(options, kgrpc.WithDiscovery(discovery), kgrpc.WithEndpoint("discovery:///"+service))
		conn, err := kgrpc.DialInsecure(ctx, options...)
		if err != nil {
			cleanup()
			return nil, nil, err
		}
		return conn, cleanup, nil
	default:
		return nil, nil, fmt.Errorf("unsupported grpc client mode: %s", mode)
	}
}
