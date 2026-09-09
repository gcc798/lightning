package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Consul struct {
	address string
	client  *http.Client
}

func NewConsul(address string) *Consul {
	return &Consul{address: strings.TrimRight(address, "/"), client: http.DefaultClient}
}

func (r *Consul) Register(ctx context.Context, instance ServiceInstance) error {
	body := map[string]any{"ID": instance.ID, "Name": instance.Name, "Meta": instance.Endpoints}
	if endpoint := instance.Endpoints[EndpointHTTP]; endpoint != "" {
		body["Check"] = map[string]any{"HTTP": strings.TrimRight(endpoint, "/") + "/health/live", "Interval": "5s", "Timeout": "2s", "DeregisterCriticalServiceAfter": "30s"}
	}
	return r.put(ctx, "/v1/agent/service/register", body)
}

func (r *Consul) Deregister(ctx context.Context, instance ServiceInstance) error {
	return r.put(ctx, "/v1/agent/service/deregister/"+instance.ID, nil)
}

func (r *Consul) Resolve(ctx context.Context, name string) ([]ServiceInstance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.address+"/v1/health/service/"+name+"?passing=true", nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("resolve %s from consul: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("resolve %s from consul: status %s", name, resp.Status)
	}
	var rows []struct {
		Service struct {
			ID   string            `json:"ID"`
			Meta map[string]string `json:"Meta"`
		} `json:"Service"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode consul response: %w", err)
	}
	instances := make([]ServiceInstance, 0, len(rows))
	for _, row := range rows {
		instances = append(instances, ServiceInstance{ID: row.Service.ID, Name: name, Endpoints: row.Service.Meta})
	}
	return instances, nil
}

func (r *Consul) Watch(_ context.Context, name string) (Watcher, error) {
	return newPollingWatcher(r, name), nil
}

func (r *Consul) Close() error { return nil }

func (r *Consul) put(ctx context.Context, path string, value any) error {
	var body bytes.Buffer
	if value != nil {
		if err := json.NewEncoder(&body).Encode(value); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, r.address+path, &body)
	if err != nil {
		return err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("consul request %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("consul request %s: status %s", path, resp.Status)
	}
	return nil
}
