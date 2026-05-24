package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qnguyenhong/automation-platform/internal/worker"
	"github.com/qnguyenhong/automation-platform/internal/worker/executors"
	"github.com/qnguyenhong/automation-platform/pkg/config"
	"github.com/qnguyenhong/automation-platform/pkg/logger"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lg := logger.New(cfg.Log.Level, cfg.Log.Format)

	registry := worker.NewRegistry()
	registry.Register(executors.NewLoadExecutor(registry, &cfg.Worker))
	registry.Register(executors.NewHTTPExecutor(30 * time.Second))

	agent := worker.NewAgent(&cfg.Worker, lg, registry)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		lg.Info("shutting down worker...")
		cancel()
	}()

	// Start worker HTTP server for server ping-pong check
	go func() {
		port := "9090"
		if envPort := os.Getenv("WORKER_PING_PORT"); envPort != "" {
			port = envPort
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"pong"}`))
		})
		lg.Info("starting worker ping server", "port", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil && err != http.ErrServerClosed {
			lg.Error("worker ping server failed", "error", err)
		}
	}()

	if err := agent.Start(ctx); err != nil {
		lg.Error("worker error", "error", err)
		os.Exit(1)
	}
}
