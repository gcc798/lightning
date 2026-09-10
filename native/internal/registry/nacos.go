package registry

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Nacos struct {
	client naming_client.INamingClient
	group  string
}

func NewNacos(options Options) (*Nacos, error) {
	servers, err := parseNacosServers(options.Address)
	if err != nil {
		return nil, err
	}
	runtimeDir := filepath.Join(os.TempDir(), "microservice-kit-nacos-"+strconv.Itoa(os.Getpid()))
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig: &constant.ClientConfig{
			NamespaceId:         options.Namespace,
			Username:            options.Username,
			Password:            options.Password,
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			DisableUseSnapShot:  true,
			CacheDir:            runtimeDir,
			LogDir:              runtimeDir,
			LogLevel:            "error",
		},
		ServerConfigs: servers,
	})
	if err != nil {
		return nil, fmt.Errorf("create nacos client: %w", err)
	}
	group := options.Group
	if group == "" {
		group = constant.DEFAULT_GROUP
	}
	return &Nacos{client: client, group: group}, nil
}

func (r *Nacos) Register(_ context.Context, instance ServiceInstance) error {
	ip, port, err := nacosInstanceAddress(instance)
	if err != nil {
		return err
	}
	metadata, err := encodeMetadata(instance)
	if err != nil {
		return err
	}
	ok, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip: ip, Port: port, ServiceName: instance.Name, GroupName: r.group,
		Weight: 1, Enable: true, Healthy: true, Ephemeral: true, Metadata: metadata,
	})
	if err != nil {
		return fmt.Errorf("register %s in nacos: %w", instance.Name, err)
	}
	if !ok {
		return fmt.Errorf("register %s in nacos: rejected", instance.Name)
	}
	return nil
}

func (r *Nacos) Deregister(_ context.Context, instance ServiceInstance) error {
	ip, port, err := nacosInstanceAddress(instance)
	if err != nil {
		return err
	}
	ok, err := r.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip: ip, Port: port, ServiceName: instance.Name, GroupName: r.group, Ephemeral: true,
	})
	if err != nil {
		return fmt.Errorf("deregister %s from nacos: %w", instance.Name, err)
	}
	if !ok {
		return fmt.Errorf("deregister %s from nacos: rejected", instance.Name)
	}
	return nil
}

func (r *Nacos) Resolve(_ context.Context, name string) ([]ServiceInstance, error) {
	rows, err := r.client.SelectInstances(vo.SelectInstancesParam{ServiceName: name, GroupName: r.group, HealthyOnly: true})
	if err != nil {
		return nil, fmt.Errorf("resolve %s from nacos: %w", name, err)
	}
	instances := make([]ServiceInstance, 0, len(rows))
	for _, row := range rows {
		id, endpoints, routes, err := decodeMetadata(row.Metadata)
		if err != nil {
			return nil, fmt.Errorf("decode nacos service %s instance %s: %w", name, row.InstanceId, err)
		}
		if id == "" {
			id = row.InstanceId
		}
		instances = append(instances, ServiceInstance{ID: id, Name: name, Endpoints: endpoints, Routes: routes})
	}
	return instances, nil
}

func (r *Nacos) Watch(_ context.Context, name string) (Watcher, error) {
	return newPollingWatcher(r, name), nil
}

func (r *Nacos) Close() error {
	r.client.CloseClient()
	return nil
}

func nacosInstanceAddress(instance ServiceInstance) (string, uint64, error) {
	endpoint := instance.Endpoints[EndpointHTTP]
	if endpoint == "" {
		endpoint = instance.Endpoints[EndpointGRPC]
		if endpoint != "" && !strings.Contains(endpoint, "://") {
			endpoint = "tcp://" + endpoint
		}
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Hostname() == "" || parsed.Port() == "" {
		return "", 0, fmt.Errorf("service %q requires an advertised HTTP or gRPC endpoint for nacos", instance.Name)
	}
	port, err := strconv.ParseUint(parsed.Port(), 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("parse nacos port for %s: %w", instance.Name, err)
	}
	return parsed.Hostname(), port, nil
}

func parseNacosServers(address string) ([]constant.ServerConfig, error) {
	parts := strings.Split(address, ",")
	servers := make([]constant.ServerConfig, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, "://") {
			part = "http://" + part
		}
		parsed, err := url.Parse(part)
		if err != nil || parsed.Hostname() == "" {
			return nil, fmt.Errorf("invalid nacos address %q", part)
		}
		port := uint64(8848)
		if parsed.Port() != "" {
			port, err = strconv.ParseUint(parsed.Port(), 10, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid nacos address %q: %w", part, err)
			}
		}
		contextPath := strings.TrimRight(parsed.Path, "/")
		if contextPath == "" {
			contextPath = "/nacos"
		}
		servers = append(servers, constant.ServerConfig{
			Scheme: parsed.Scheme, IpAddr: parsed.Hostname(), Port: port, ContextPath: contextPath,
		})
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("nacos address is required")
	}
	return servers, nil
}
