package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	resourcev1 "github.com/gcc798/lightning/internal/api/resource/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	_ "github.com/gcc798/lightning/internal/openapi"
	"github.com/gcc798/lightning/internal/registry"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type gateway struct {
	registry registry.Registry
	selector *registry.Selector
	security iamv1.API
}

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
	service := routeService(request.URL.Path)
	if service == "" {
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
		request.Header.Set("X-Lightning-User-ID", strconv.FormatInt(claims.UserId, 10))
	}
	instances, err := g.registry.Resolve(request.Context(), service)
	if err != nil {
		writeError(writer, http.StatusBadGateway, "服务发现失败")
		return
	}
	instance, err := g.selector.Pick(service, instances)
	if err != nil {
		writeError(writer, http.StatusServiceUnavailable, "服务暂不可用")
		return
	}
	target, err := url.Parse(instance.Endpoints[registry.EndpointHTTP])
	if err != nil || target.Scheme == "" || target.Host == "" {
		writeError(writer, http.StatusBadGateway, "服务地址无效")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = otelhttp.NewTransport(http.DefaultTransport)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
		writeError(w, http.StatusBadGateway, "服务调用失败")
	}
	proxy.ServeHTTP(writer, request)
}

func routeService(path string) string {
	switch {
	case path == "/login", path == "/logout", strings.HasPrefix(path, "/auth/"), strings.HasPrefix(path, "/captcha/"),
		path == "/metrics",
		strings.HasPrefix(path, "/api/v1/user"), strings.HasPrefix(path, "/api/v1/role"),
		strings.HasPrefix(path, "/api/v1/menu"), strings.HasPrefix(path, "/api/v1/org"),
		strings.HasPrefix(path, "/api/v1/api-permission"), strings.HasPrefix(path, "/system/user"),
		path == "/resource/websocket", path == "/resource/sms/code":
		return iamv1.ServiceName
	case strings.HasPrefix(path, "/api/v1/config"), strings.HasPrefix(path, "/api/v1/dict"),
		strings.HasPrefix(path, "/api/v1/loginLog"), strings.HasPrefix(path, "/api/v1/operLog"):
		return sysv1.ServiceName
	case strings.HasPrefix(path, "/api/v1/attachment"):
		return resourcev1.ServiceName
	default:
		return ""
	}
}

func publicPath(path string) bool {
	return path == "/login" || path == "/auth/refresh" || path == "/resource/sms/code" || strings.HasPrefix(path, "/captcha/")
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]any{"code": status, "msg": message})
}
