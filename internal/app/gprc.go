package app

import (
	"context"
	grpcController "ton-node/internal/controller/grpc"
	"ton-node/internal/controller/grpc/tonnodepb"
	"ton-node/internal/domain/logger"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(dep *App, l logger.Loggerer) (*grpc.Server, *prometheus.Registry) {
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
			grpcController.AuthMiddleware(),
			grpcController.LoggerInterceptor(l),
		),
	)

	server := grpc.NewServer(opts...)
	tonnodepb.RegisterTonNodeServiceServer(server, dep.Controller)

	reflection.Register(server)

	return server, reg
}
