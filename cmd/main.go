package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dmdutra/github-outage-exporter/internal/exporter"
	"github.com/dmdutra/github-outage-exporter/internal/githubstatus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenAddr := envOrDefault("LISTEN_ADDR", ":8080")
	githubStatusURL := envOrDefault("GITHUB_STATUS_URL", "https://www.githubstatus.com")
	scrapeTimeout := envDurationOrDefault("SCRAPE_TIMEOUT", 10*time.Second)
	metricsPath := envOrDefault("METRICS_PATH", "/metrics")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	client := githubstatus.NewClient(githubStatusURL, scrapeTimeout)
	collector := exporter.NewCollector(client, logger)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collector,
	)

	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, metricsPath, http.StatusFound)
	})

	logger.Info("starting github outage exporter",
		"listen", listenAddr,
		"github_status_url", githubStatusURL,
		"metrics_path", metricsPath,
		"scrape_timeout", scrapeTimeout,
	)

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
