package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const (
	EndpointHTTP = "http"
	EndpointGRPC = "grpc"
)

type ServiceInstance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Endpoints map[string]string `json:"endpoints"`
}

type Watcher interface {
	Next(context.Context) ([]ServiceInstance, error)
	Stop()
}

type Registry interface {
	Register(context.Context, ServiceInstance) error
	Deregister(context.Context, ServiceInstance) error
	Resolve(context.Context, string) ([]ServiceInstance, error)
	Watch(context.Context, string) (Watcher, error)
	Close() error
}

func New(driver, address, prefix string) (Registry, error) {
	switch driver {
	case "consul":
		return NewConsul(address), nil
	case "etcd":
		return NewEtcd(address, prefix), nil
	case "inprocess", "":
		return NewInProcess(), nil
	default:
		return nil, fmt.Errorf("unsupported registry driver %q", driver)
	}
}

type Selector struct {
	mu   sync.Mutex
	next map[string]uint64
}

func NewSelector() *Selector { return &Selector{next: make(map[string]uint64)} }

func (s *Selector) Pick(name string, instances []ServiceInstance) (ServiceInstance, error) {
	if len(instances) == 0 {
		return ServiceInstance{}, fmt.Errorf("service %q has no healthy instances", name)
	}
	s.mu.Lock()
	index := s.next[name] % uint64(len(instances))
	s.next[name]++
	s.mu.Unlock()
	return instances[index], nil
}

var ErrWatcherStopped = errors.New("registry watcher stopped")
