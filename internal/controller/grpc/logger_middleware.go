package grpc

import (
	"context"
	"log/slog"
	"time"
	"ton-node/internal/domain/logger"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggerInterceptor(log logger.Loggerer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		method := info.FullMethod

		requestID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids, exists := md["X-Request-Id"]; exists && len(ids) > 0 {
				requestID = ids[0]
			}
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		logWithCtx := log.With(slog.Group("request_info",
			slog.String("id", requestID),
			slog.String("method", method),
			slog.Any("body", req),
		))

		resp, err := handler(ctx, req)

		logWithCtx = logWithCtx.With(
			slog.Group("response_info",
				slog.Any("body", resp),
				slog.Int64("latency", time.Since(start).Milliseconds()),
			),
		)

		if err != nil {
			logWithCtx.Error("request failed", slog.String("error", err.Error()))
		} else {
			logWithCtx.Info("request")
		}

		return resp, err
	}
}
