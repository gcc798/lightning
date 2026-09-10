package main

import (
	"flag"
	"os"

	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/pkg/configx"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	configPath := flag.String("conf", "application/sys-rpc/configs/config.yaml", "config path")
	flag.Parse()
	cfg := mustLoadConfig(*configPath)
	logger := log.With(
		log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.name", "sys-rpc",
	)
	app, appCleanup, err := wireApp(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer appCleanup()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func mustLoadConfig(path string) *conf.Bootstrap {
	var cfg conf.Bootstrap
	configx.MustLoad(path, &cfg)
	return &cfg
}
