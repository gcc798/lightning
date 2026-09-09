package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	sysdomain "github.com/gcc798/lightning/application/sys/internal/domain"
	"github.com/gcc798/lightning/application/sys/internal/migrations"
	"github.com/gcc798/lightning/application/sys/internal/router"
	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/httpserver"
	"github.com/gcc798/lightning/internal/httpx"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/registry"
	"github.com/gcc798/lightning/internal/transport"
	"google.golang.org/grpc"
)

func main() {
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()
	cfg, v, err := config.Load("application/sys", config.ServiceSystem)
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
	cont, err := container.New(cfg, v, log, container.WithSystemInfrastructure())
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
	reg, err := registry.New(cfg.Registry.Driver, cfg.Registry.Address, cfg.Registry.Prefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	pool := transport.NewClientPool(reg)
	defer pool.Close()
	defer reg.Close()
	security := iamv1.NewCached(iamv1.NewRemote(pool), 5*time.Second)
	systemAPI := sysdomain.NewAPI(cont.GetDB(), cont.GetLogger())
	grpcServer, err := transport.StartRegisteredGRPC(ctx, cfg, sysv1.ServiceName, func(server *grpc.Server) {
		sysv1.RegisterSystemServiceServer(server, sysdomain.NewGRPCServer(systemAPI))
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
	if err := httpserver.RunHTTP(ctx, cont, security, systemAPI, func(r *httpx.Router) error { return router.Setup(r, cont, security, systemAPI) }); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
