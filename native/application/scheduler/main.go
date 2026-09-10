package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcc798/microservice-kit/application/scheduler/jobs"
	resourcev1 "github.com/gcc798/microservice-kit/internal/api/resource/v1"
	sysv1 "github.com/gcc798/microservice-kit/internal/api/sys/v1"
	"github.com/gcc798/microservice-kit/internal/config"
	"github.com/gcc798/microservice-kit/internal/container"
	logging "github.com/gcc798/microservice-kit/internal/logger"
	"github.com/gcc798/microservice-kit/internal/modules"
	"github.com/gcc798/microservice-kit/internal/registry"
	"github.com/gcc798/microservice-kit/internal/telemetry"
	"github.com/gcc798/microservice-kit/internal/transport"
	"go.uber.org/zap"
)

func main() {
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()
	cfg, v, err := config.Load("application/scheduler", config.ServiceScheduler)
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
	shutdownTelemetry, err := telemetry.Init(context.Background(), string(config.ServiceScheduler), cfg.Service.ID, config.CurrentEnv())
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
	cont, err := container.New(cfg, v, log)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	reg, err := registry.New(registry.Options{
		Driver: cfg.Registry.Driver, Address: cfg.Registry.Address, Prefix: cfg.Registry.Prefix,
		Namespace: cfg.Registry.Namespace, Group: cfg.Registry.Group, Username: cfg.Registry.Username, Password: cfg.Registry.Password,
	})
	if err != nil {
		log.Error("scheduler process exited with error", zap.Error(err))
		exitCode = 1
		return
	}
	defer reg.Close()
	pool := transport.NewClientPool(reg)
	defer pool.Close()
	definitions := jobs.Definitions(sysv1.NewRemote(pool), resourcev1.NewRemote(pool), cont.GetLogger())
	if err := cont.RegisterModules(ctx,
		modules.NewSchedulerModule(definitions),
	); err != nil {
		log.Error("scheduler process exited with error", zap.Error(fmt.Errorf("initialize scheduler modules: %w", err)))
		exitCode = 1
		return
	}
	if err := cont.StartModules(ctx); err != nil {
		log.Error("scheduler process exited with error", zap.Error(fmt.Errorf("start scheduler modules: %w", err)))
		exitCode = 1
		return
	}
	cont.GetLogger().Info("scheduler process started")
	<-ctx.Done()
	cont.GetLogger().Info("stopping scheduler process")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := cont.StopModules(shutdownCtx); err != nil {
		log.Error("scheduler process exited with error", zap.Error(err))
		exitCode = 1
		return
	}
}
