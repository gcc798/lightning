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
	Routes    []HTTPRoute       `json:"routes,omitempty"`
}

type HTTPRoute struct {
	Method string `json:"m"`
	Path   string `json:"p"`
}

type Options struct {
	Driver    string
	Address   string
	Prefix    string
	Namespace string
	Group     string
	Username  string
	Password  string
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

func New(options Options) (Registry, error) {
	switch options.Driver {
	case "consul":
		return NewConsul(options.Address), nil
	case "etcd":
		return NewEtcd(options.Address, options.Prefix), nil
	case "nacos":
		return NewNacos(options)
	case "inprocess", "":
		// Unit tests only: each process gets an isolated registry.
		return NewInProcess(), nil
	default:
		return nil, fmt.Errorf("unsupported registry driver %q", options.Driver)
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
