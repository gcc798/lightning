// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gcc798/lightning/application/sys-api/internal/config"
	"github.com/gcc798/lightning/application/sys-api/internal/handler"
	"github.com/gcc798/lightning/application/sys-api/internal/svc"
	"github.com/gcc798/lightning/common/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/sys-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		middleware.SetCORSHeaders(header, middleware.CORSConfig{})
	}, nil, "*"))
	defer server.Stop()

	server.Use(middleware.PanicRecoveryMiddleware)
	server.Use(middleware.NewJWTAuthMiddleware(middleware.JWTAuthConfig{
		Secret:      c.Jwt.Secret,
		TokenHeader: c.Auth.TokenHeader,
		WhiteList: []string{
			"/login",
			"/logout",
			"/auth/login",
			"/auth/logout",
			"/auth/refresh",
			"/captcha/*",
			"/resource/sms/code",
			"/health",
			"/health/ready",
			"/health/live",
			"/health/startup",
		},
	}).Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
