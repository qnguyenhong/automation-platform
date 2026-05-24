package main

import (
	"context"
	"log"
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

	if err := agent.Start(ctx); err != nil {
		lg.Error("worker error", "error", err)
		os.Exit(1)
	}
}
