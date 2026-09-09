package transport

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/gcc798/lightning/internal/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(address string) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen for grpc: %w", err)
	}
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	return &GRPCServer{server: server, listener: listener}, nil
}

func (s *GRPCServer) Server() *grpc.Server { return s.server }
func (s *GRPCServer) Address() string      { return s.listener.Addr().String() }
func (s *GRPCServer) Serve() error         { return s.server.Serve(s.listener) }
func (s *GRPCServer) GracefulStop()        { s.server.GracefulStop() }

type ClientPool struct {
	registry registry.Registry
	selector *registry.Selector
	mu       sync.Mutex
	conns    map[string]*grpc.ClientConn
}

func NewClientPool(reg registry.Registry) *ClientPool {
	return &ClientPool{registry: reg, selector: registry.NewSelector(), conns: make(map[string]*grpc.ClientConn)}
}

func (p *ClientPool) Conn(ctx context.Context, service string) (*grpc.ClientConn, error) {
	instances, err := p.registry.Resolve(ctx, service)
	if err != nil {
		return nil, err
	}
	instance, err := p.selector.Pick(service, instances)
	if err != nil {
		return nil, err
	}
	address := instance.Endpoints[registry.EndpointGRPC]
	if address == "" {
		return nil, fmt.Errorf("service %q instance %q has no grpc endpoint", service, instance.ID)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if conn := p.conns[address]; conn != nil {
		return conn, nil
	}
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to %s at %s: %w", service, address, err)
	}
	p.conns[address] = conn
	return conn, nil
}

func (p *ClientPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	var first error
	for address, conn := range p.conns {
		if err := conn.Close(); err != nil && first == nil {
			first = err
		}
		delete(p.conns, address)
	}
	return first
}
