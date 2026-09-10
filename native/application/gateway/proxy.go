package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	iamv1 "github.com/gcc798/microservice-kit/internal/api/iam/v1"
	resourcev1 "github.com/gcc798/microservice-kit/internal/api/resource/v1"
	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	_ "github.com/gcc798/microservice-kit/internal/openapi"
	"github.com/gcc798/microservice-kit/internal/registry"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type gateway struct {
	registry registry.Registry
	selector *registry.Selector
	security iamv1.API
	routes   atomic.Pointer[routeTable]
}

type routeTable struct {
	entries []routeTarget
}

type routeTarget struct {
	method    string
	pattern   string
	service   string
	instances []registry.ServiceInstance
}

var routedServices = []string{iamv1.ServiceName, sysv1.ServiceName, resourcev1.ServiceName}

func (g *gateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/health" || strings.HasPrefix(request.URL.Path, "/health/") {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"status":"healthy"}`))
		return
	}
	if strings.HasPrefix(request.URL.Path, "/swagger/") {
		httpSwagger.Handler().ServeHTTP(writer, request)
		return
	}
	target, err := g.route(request.Method, request.URL.Path)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	if !publicPath(request.URL.Path) {
		token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
		if token == "" && request.URL.Path == "/resource/websocket" {
			token = request.URL.Query().Get("Authorization")
		}
		claims, err := g.security.ValidateAccessToken(request.Context(), token)
		if err != nil {
			writeError(writer, http.StatusUnauthorized, "未登录或登录已失效")
			return
		}
		request.Header.Set("X-Microservice-Kit-User-ID", strconv.FormatInt(claims.UserId, 10))
	}
	instance, err := g.selector.Pick(target.service, target.instances)
	if err != nil {
		writeError(writer, http.StatusServiceUnavailable, "服务暂不可用")
		return
	}
	targetURL, err := url.Parse(instance.Endpoints[registry.EndpointHTTP])
	if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
		writeError(writer, http.StatusBadGateway, "服务地址无效")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = otelhttp.NewTransport(http.DefaultTransport)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
		writeError(w, http.StatusBadGateway, "服务调用失败")
	}
	proxy.ServeHTTP(writer, request)
}

func (g *gateway) refreshRoutes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	services := make(map[string][]registry.ServiceInstance, len(routedServices))
	for _, service := range routedServices {
		instances, err := g.registry.Resolve(ctx, service)
		if err != nil {
			return err
		}
		services[service] = instances
	}
	table, err := buildRouteTable(services)
	if err != nil {
		return err
	}
	g.routes.Store(table)
	return nil
}

func buildRouteTable(services map[string][]registry.ServiceInstance) (*routeTable, error) {
	byRoute := make(map[string]*routeTarget)
	for service, instances := range services {
		for _, instance := range instances {
			for _, route := range instance.Routes {
				if route.Method == "" || route.Path == "" || route.Path == "/metrics" || route.Path == "/health" || strings.HasPrefix(route.Path, "/health/") {
					continue
				}
				key := route.Method + "\x00" + canonicalRoutePath(route.Path)
				target := byRoute[key]
				if target == nil {
					target = &routeTarget{method: route.Method, pattern: route.Path, service: service}
					byRoute[key] = target
				} else if target.service != service {
					return nil, fmt.Errorf("HTTP route %s %s is registered by both %s and %s", route.Method, route.Path, target.service, service)
				}
				if !hasInstance(target.instances, instance.ID) {
					target.instances = append(target.instances, instance)
				}
			}
		}
	}
	entries := make([]routeTarget, 0, len(byRoute))
	for _, target := range byRoute {
		entries = append(entries, *target)
	}
	sort.Slice(entries, func(i, j int) bool {
		return routeLess(entries[i].pattern, entries[j].pattern)
	})
	return &routeTable{entries: entries}, nil
}

func (g *gateway) route(method, path string) (routeTarget, error) {
	table := g.routes.Load()
	if table != nil {
		for _, target := range table.entries {
			if target.method == method && matchRoutePath(target.pattern, path) {
				return target, nil
			}
		}
	}
	return routeTarget{}, fmt.Errorf("route not found")
}

func matchRoutePath(pattern, path string) bool {
	patterns := strings.Split(strings.TrimPrefix(pattern, "/"), "/")
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for index, segment := range patterns {
		if segment == "*" {
			return true
		}
		if index >= len(parts) || parts[index] == "" {
			return false
		}
		if segment != parts[index] && !strings.HasPrefix(segment, ":") {
			return false
		}
	}
	return len(patterns) == len(parts)
}

func routeLess(left, right string) bool {
	leftParts := strings.Split(strings.TrimPrefix(left, "/"), "/")
	rightParts := strings.Split(strings.TrimPrefix(right, "/"), "/")
	for index := 0; index < min(len(leftParts), len(rightParts)); index++ {
		leftRank := routeSegmentRank(leftParts[index])
		rightRank := routeSegmentRank(rightParts[index])
		if leftRank != rightRank {
			return leftRank > rightRank
		}
	}
	if len(leftParts) != len(rightParts) {
		return len(leftParts) > len(rightParts)
	}
	return left < right
}

func routeSegmentRank(segment string) int {
	if segment == "*" {
		return 0
	}
	if strings.HasPrefix(segment, ":") {
		return 1
	}
	return 2
}

func hasInstance(instances []registry.ServiceInstance, id string) bool {
	for _, instance := range instances {
		if instance.ID == id {
			return true
		}
	}
	return false
}

func canonicalRoutePath(pattern string) string {
	parts := strings.Split(pattern, "/")
	for index, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[index] = ":"
		}
	}
	return strings.Join(parts, "/")
}

func publicPath(path string) bool {
	return path == "/login" || path == "/auth/refresh" || path == "/resource/sms/code" || strings.HasPrefix(path, "/captcha/")
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]any{"code": status, "msg": message})
}
