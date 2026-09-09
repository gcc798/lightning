package transport

import (
	"context"
	"errors"
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
	host := cfg.Service.AdvertiseHost
	go func() { _ = server.Serve() }()
	instance, err := RegisterService(ctx, reg, name, cfg.Service.ID, map[string]string{registry.EndpointHTTP: "http://" + host + ":" + strconv.Itoa(cfg.Server.Port), registry.EndpointGRPC: host + ":" + strconv.Itoa(cfg.GRPC.Port)})
	if err != nil {
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
