package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ton-node/config"
	log "ton-node/infra/logger"
	"ton-node/infra/node"
	grpcController "ton-node/internal/controller/grpc"
	"ton-node/internal/controller/grpc/tonnodepb"
	"ton-node/internal/domain/logger"
	tonModel "ton-node/internal/domain/ton"
	"ton-node/internal/domain/usecase"
	uc "ton-node/internal/usecase"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc/reflection"

	"emperror.dev/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type App struct {
	Controller *grpcController.TonNodeController
	UseCase    usecase.UseCase
	NodeSvc    tonModel.NodeService
}

func NewApp(nodeSvc tonModel.NodeService, l logger.Loggerer) *App {
	useCase := uc.NewUseCase(nodeSvc)
	controller := grpcController.NewTonNodeController(useCase, l)

	return &App{
		NodeSvc:    nodeSvc,
		UseCase:    useCase,
		Controller: controller,
	}
}

func NewNodeService(rpcURL, apiKey string, timeout time.Duration, rps, burst int) tonModel.NodeService {
	return node.NewService(rpcURL, apiKey, timeout, rps, burst)
}

func NewGRPCServer(dep *App, l logger.Loggerer) *grpc.Server {
	var opts []grpc.ServerOption

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	exemplarFromContext := func(ctx context.Context) prometheus.Labels {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return prometheus.Labels{"traceID": span.TraceID().String()}
		}
		return nil
	}

	opts = append(opts,
		grpc.ChainUnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(exemplarFromContext)),
			// grpcController.AuthMiddleware(),
			grpcController.LoggerInterceptor(l),
		),
	)

	grpcServer := grpc.NewServer(opts...)
	tonnodepb.RegisterTonNodeServiceServer(grpcServer, dep.Controller)

	reflection.Register(grpcServer)

	return grpcServer
}

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

	nodeSvc := NewNodeService(
		viper.GetString("app.node.url"),
		viper.GetString("app.node.api_key"),
		viper.GetDuration("app.node.timeout"),
		viper.GetInt("app.node.rate_limit"),
		viper.GetInt("app.node.rate_burst"),
	)

	di := NewApp(nodeSvc, l)

	addr := fmt.Sprintf("%s:%s",
		viper.GetString("app.server.host"),
		viper.GetString("app.server.port"),
	)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error("Failed to listen", err)
		os.Exit(1)
	}

	grpcServer := NewGRPCServer(di, l)

	l.Info("starting server", slog.String("address", addr))

	shutdownTimeout := viper.GetDuration("app.server.shutdown_timeout")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		if err = grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serveErr <- err
		} else {
			serveErr <- nil
		}
	}()

	select {
	case <-ctx.Done():
		l.Info("shutdown signal received")
	case err = <-serveErr:
		if err != nil {
			l.Error("gRPC server failed", slog.Any("error", err))
		}
		return
	}

	l.Info("shutdown initiated")

	gracefulCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	stopped := make(chan struct{})

	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		l.Info("server stopped gracefully")
	case <-gracefulCtx.Done():
		l.Warn("shutdown timeout, forcing stop")

		select {
		case <-stopped:
			l.Info("server stopped just after timeout")
		default:
			grpcServer.Stop()
		}
	}
}
