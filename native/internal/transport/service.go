package transport

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/registry"
	"google.golang.org/grpc"
)

type RegisteredGRPC struct {
	Server   *GRPCServer
	Registry registry.Registry
	Instance registry.ServiceInstance
}

func StartRegisteredGRPC(ctx context.Context, cfg *config.Config, name string, register func(*grpc.Server)) (*RegisteredGRPC, error) {
	reg, err := registry.New(cfg.Registry.Driver, cfg.Registry.Address, cfg.Registry.Prefix)
	if err != nil {
		return nil, err
	}
	server, err := NewGRPCServer(":" + strconv.Itoa(cfg.GRPC.Port))
	if err != nil {
		_ = reg.Close()
		return nil, err
	}
	register(server.Server())
	id := cfg.Service.ID
	if id == "" {
		host, _ := os.Hostname()
		id = fmt.Sprintf("%s-%s-%d", name, host, cfg.Server.Port)
	}
	host := cfg.Service.AdvertiseHost
	instance := registry.ServiceInstance{ID: id, Name: name, Endpoints: map[string]string{registry.EndpointHTTP: "http://" + host + ":" + strconv.Itoa(cfg.Server.Port), registry.EndpointGRPC: host + ":" + strconv.Itoa(cfg.GRPC.Port)}}
	go func() { _ = server.Serve() }()
	if err := reg.Register(ctx, instance); err != nil {
		server.GracefulStop()
		_ = reg.Close()
		return nil, err
	}
	return &RegisteredGRPC{Server: server, Registry: reg, Instance: instance}, nil
}
func (s *RegisteredGRPC) Stop(ctx context.Context) error {
	err := s.Registry.Deregister(ctx, s.Instance)
	s.Server.GracefulStop()
	if e := s.Registry.Close(); err == nil {
		err = e
	}
	return err
}
