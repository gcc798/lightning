package transport

import (
	"context"
	"errors"
	"strconv"

	"github.com/gcc798/microservice-kit/internal/config"
	"github.com/gcc798/microservice-kit/internal/registry"
	"google.golang.org/grpc"
)

type RegisteredGRPC struct {
	Server   *GRPCServer
	Registry registry.Registry
	Instance registry.ServiceInstance
}

func RegisterService(ctx context.Context, reg registry.Registry, name, id string, endpoints map[string]string) (registry.ServiceInstance, error) {
	if id == "" {
		return registry.ServiceInstance{}, errors.New("service instance ID is required")
	}
	if endpoints == nil {
		endpoints = map[string]string{}
	}
	instance := registry.ServiceInstance{ID: id, Name: name, Endpoints: endpoints}
	return instance, reg.Register(ctx, instance)
}

func StartRegisteredGRPC(ctx context.Context, reg registry.Registry, cfg *config.Config, name string, routes []registry.HTTPRoute, register func(*grpc.Server)) (*RegisteredGRPC, error) {
	server, err := NewGRPCServer(":" + strconv.Itoa(cfg.GRPC.Port))
	if err != nil {
		return nil, err
	}
	register(server.Server())
	host := cfg.Service.AdvertiseHost
	go func() { _ = server.Serve() }()
	instance := registry.ServiceInstance{
		ID: cfg.Service.ID, Name: name, Routes: routes,
		Endpoints: map[string]string{registry.EndpointHTTP: "http://" + host + ":" + strconv.Itoa(cfg.Server.Port), registry.EndpointGRPC: host + ":" + strconv.Itoa(cfg.GRPC.Port)},
	}
	if err := reg.Register(ctx, instance); err != nil {
		server.GracefulStop()
		return nil, err
	}
	return &RegisteredGRPC{Server: server, Registry: reg, Instance: instance}, nil
}
func (s *RegisteredGRPC) Stop(ctx context.Context) error {
	err := s.Registry.Deregister(ctx, s.Instance)
	s.Server.GracefulStop()
	return err
}
