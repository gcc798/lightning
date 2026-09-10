package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	resourcedomain "github.com/gcc798/microservice-kit/application/resource/internal/domain"
	"github.com/gcc798/microservice-kit/application/resource/internal/migrations"
	"github.com/gcc798/microservice-kit/application/resource/internal/router"
	iamv1 "github.com/gcc798/microservice-kit/internal/api/iam/v1"
	resourcev1 "github.com/gcc798/microservice-kit/internal/api/resource/v1"
	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	"github.com/gcc798/microservice-kit/internal/config"
	"github.com/gcc798/microservice-kit/internal/container"
	"github.com/gcc798/microservice-kit/internal/httpserver"
	"github.com/gcc798/microservice-kit/internal/httpx"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/registry"
	"github.com/gcc798/microservice-kit/internal/telemetry"
	"github.com/gcc798/microservice-kit/internal/transport"
	"google.golang.org/grpc"
)

func main() {
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()
	cfg, v, err := config.Load("application/resource", config.ServiceResource)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	log, err := logging.NewLogger(config.CurrentEnv(), cfg.AppDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	shutdownTelemetry, err := telemetry.Init(context.Background(), string(config.ServiceResource), cfg.Service.ID, config.CurrentEnv())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	defer func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTelemetry(shutdown)
	}()
	cont, err := container.New(cfg, v, log, container.WithResourceInfrastructure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	sqlDB, err := cont.GetDB().DB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	if err := migrations.Up(sqlDB); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	reg, err := registry.New(registry.Options{
		Driver: cfg.Registry.Driver, Address: cfg.Registry.Address, Prefix: cfg.Registry.Prefix,
		Namespace: cfg.Registry.Namespace, Group: cfg.Registry.Group, Username: cfg.Registry.Username, Password: cfg.Registry.Password,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	pool := transport.NewClientPool(reg)
	defer pool.Close()
	defer reg.Close()
	security := iamv1.NewCached(iamv1.NewRemote(pool), 5*time.Second)
	systemAPI := sysv1.NewRemote(pool)
	attachments := resourcedomain.NewAttachmentService(cont.GetDB(), cont.GetStorage(), cont.GetLogger())
	resourceAPI := resourcedomain.NewAPI(attachments)
	httpServer, routes, err := httpserver.New(cont, systemAPI, func(r *httpx.Router) error { return router.Setup(r, cont, security, systemAPI) })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	grpcServer, err := transport.StartRegisteredGRPC(ctx, reg, cfg, resourcev1.ServiceName, routes, func(server *grpc.Server) {
		resourcev1.RegisterResourceServiceServer(server, resourcedomain.NewGRPCServer(resourceAPI))
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	defer func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = grpcServer.Stop(shutdown)
	}()
	if err := httpServer.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
