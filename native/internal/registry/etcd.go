package registry

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Etcd struct {
	address string
	prefix  string
	client  *http.Client
	mu      sync.Mutex
	leases  map[string]*etcdLease
}

type etcdLease struct {
	id     string
	cancel context.CancelFunc
}

func NewEtcd(address, prefix string) *Etcd {
	if prefix == "" {
		prefix = "/lightning/services"
	}
	return &Etcd{
		address: strings.TrimRight(address, "/"),
		prefix:  strings.TrimRight(prefix, "/"),
		client:  http.DefaultClient,
		leases:  make(map[string]*etcdLease),
	}
}

func (r *Etcd) Register(ctx context.Context, instance ServiceInstance) error {
	var leaseResponse struct {
		ID string `json:"ID"`
	}
	if err := r.post(ctx, "/v3/lease/grant", map[string]any{"TTL": 15}, &leaseResponse); err != nil {
		return fmt.Errorf("grant etcd lease: %w", err)
	}
	value, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	key := r.key(instance)
	if err := r.post(ctx, "/v3/kv/put", map[string]any{
		"key": base64.StdEncoding.EncodeToString([]byte(key)), "value": base64.StdEncoding.EncodeToString(value), "lease": leaseResponse.ID,
	}, nil); err != nil {
		return fmt.Errorf("register service in etcd: %w", err)
	}
	leaseCtx, cancel := context.WithCancel(context.Background())
	r.mu.Lock()
	if previous := r.leases[key]; previous != nil {
		previous.cancel()
	}
	r.leases[key] = &etcdLease{id: leaseResponse.ID, cancel: cancel}
	r.mu.Unlock()
	go r.keepAlive(leaseCtx, leaseResponse.ID)
	return nil
}

func (r *Etcd) Deregister(ctx context.Context, instance ServiceInstance) error {
	key := r.key(instance)
	r.mu.Lock()
	lease := r.leases[key]
	delete(r.leases, key)
	r.mu.Unlock()
	if lease != nil {
		lease.cancel()
	}
	return r.post(ctx, "/v3/kv/deleterange", map[string]string{"key": base64.StdEncoding.EncodeToString([]byte(key))}, nil)
}

func (r *Etcd) Resolve(ctx context.Context, name string) ([]ServiceInstance, error) {
	prefix := r.prefix + "/" + name + "/"
	end := prefixEnd(prefix)
	var response struct {
		KVs []struct {
			Value string `json:"value"`
		} `json:"kvs"`
	}
	if err := r.post(ctx, "/v3/kv/range", map[string]string{
		"key": base64.StdEncoding.EncodeToString([]byte(prefix)), "range_end": base64.StdEncoding.EncodeToString([]byte(end)),
	}, &response); err != nil {
		return nil, fmt.Errorf("resolve %s from etcd: %w", name, err)
	}
	instances := make([]ServiceInstance, 0, len(response.KVs))
	for _, row := range response.KVs {
		value, err := base64.StdEncoding.DecodeString(row.Value)
		if err != nil {
			return nil, fmt.Errorf("decode etcd service value: %w", err)
		}
		var instance ServiceInstance
		if err := json.Unmarshal(value, &instance); err != nil {
			return nil, fmt.Errorf("decode etcd service instance: %w", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (r *Etcd) Watch(_ context.Context, name string) (Watcher, error) {
	return newPollingWatcher(r, name), nil
}

func (r *Etcd) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, lease := range r.leases {
		lease.cancel()
		delete(r.leases, key)
	}
	return nil
}

func (r *Etcd) keepAlive(ctx context.Context, id string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			keepAliveCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			_ = r.post(keepAliveCtx, "/v3/lease/keepalive", map[string]any{"ID": id}, nil)
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (r *Etcd) key(instance ServiceInstance) string {
	return r.prefix + "/" + instance.Name + "/" + instance.ID
}

func (r *Etcd) post(ctx context.Context, path string, input, output any) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(input); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.address+path, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status %s", resp.Status)
	}
	if output != nil {
		return json.NewDecoder(resp.Body).Decode(output)
	}
	return nil
}

func prefixEnd(prefix string) string {
	value := []byte(prefix)
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] < 0xff {
			value[i]++
			return string(value[:i+1])
		}
	}
	return "\x00"
}
