package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"ton-node/config"
	log "ton-node/infra/logger"

	"emperror.dev/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func StartApp() {
	l := log.Init(
		log.Options{
			Level:   slog.LevelInfo,
			AppName: "ton-node",
		},
	)

	if err := config.LoadConfig(); err != nil {
		l.Error("Failed to load config", err)
		os.Exit(1)
	}

	di := newApp(l)

	addr := fmt.Sprintf("%s:%s",
		viper.GetString("app.server.host"),
		viper.GetString("app.server.port"),
	)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error("Failed to listen", err)
		os.Exit(1)
	}

	shutdownTimeout := viper.GetDuration("app.server.shutdown_timeout")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)

	server, promRegistry := NewGRPCServer(di, l)
	go func() {
		l.Info("starting server", slog.String("address", addr))
		if err = server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErr <- err
		} else {
			serveErr <- nil
		}
	}()

	metricsServer := startMetricsServer(promRegistry, l)

	select {
	case <-ctx.Done():
		l.Info("shutdown signal received")
	case err = <-serveErr:
		if err != nil {
			l.Error("gRPC server failed", slog.Any("error", err))
		}
		return
	}

	gracefulCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	grpcDone := make(chan struct{})

	go func() {
		server.GracefulStop()
		close(grpcDone)
	}()

	metricsDone := make(chan struct{})

	go func() {
		if err = metricsServer.Shutdown(gracefulCtx); err != nil {
			l.Warn("failed to shutdown metrics server", err)
		}
		close(metricsDone)
	}()

	select {
	case <-grpcDone:
	case <-gracefulCtx.Done():
		l.Warn("shutdown timeout, forcing stop")

		select {
		case <-grpcDone:
		default:
			server.Stop()
		}
	}

	<-metricsDone

	l.Info("all services stopped gracefully")
	l.Info("exit")
}
