package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcc798/lightning/application/iam/internal/domain"
	"github.com/gcc798/lightning/application/iam/internal/migrations"
	"github.com/gcc798/lightning/application/iam/internal/router"
	iamv1 "github.com/gcc798/lightning/internal/api/iam/v1"
	sysv1 "github.com/gcc798/lightning/internal/api/sys/v1"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/httpserver"
	"github.com/gcc798/lightning/internal/httpx"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gcc798/lightning/internal/modules"
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
	cfg, v, err := config.Load("application/iam", config.ServiceIAM)
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
	cont, err := container.New(cfg, v, log, container.WithIAMInfrastructure())
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
	if err := cont.RegisterModules(ctx, modules.NewSMSModule(), modules.NewEmailModule(), modules.NewWeChatModule(), modules.NewCaptchaModule(), modules.NewWebSocketModule(cfg.WebSocket.Enabled)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	if err := cont.StartModules(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	defer func() {
		shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = cont.StopModules(shutdown)
	}()

	reg, err := registry.New(cfg.Registry.Driver, cfg.Registry.Address, cfg.Registry.Prefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	pool := transport.NewClientPool(reg)
	defer pool.Close()
	defer reg.Close()
	systemAPI := sysv1.NewRemote(pool)
	tokens := iam.NewTokenManager(cont.GetJWT(), cont.GetRedis(), cont.GetLogger())
	security := iam.NewAuthzAPI(tokens, iam.NewPermissionService(cont.GetDB(), cont.GetLogger()))
	grpcServer, err := transport.StartRegisteredGRPC(ctx, cfg, iamv1.ServiceName, func(server *grpc.Server) {
		iamv1.RegisterIAMServiceServer(server, iam.NewGRPCServer(security))
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
