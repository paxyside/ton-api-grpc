package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"ton-node/internal/domain/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func startMetricsServer(registry *prometheus.Registry, l logger.Loggerer) *http.Server {
	addr := fmt.Sprintf("%s:%s",
		viper.GetString("app.server.prometheus_host"),
		viper.GetString("app.server.prometheus_port"),
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: viper.GetDuration("app.server.read_header_timeout"),
	}

	go func() {
		l.Info("metrics server listening", slog.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("metrics server failed", err)
		}
	}()

	return srv
}
