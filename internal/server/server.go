package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/qnguyenhong/automation-platform/internal/handler"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/internal/ws"
	"github.com/qnguyenhong/automation-platform/pkg/config"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	pool       *pgxpool.Pool
	hub        *ws.Hub
}

type Dependencies struct {
	Config       *config.Config
	Logger       *slog.Logger
	Pool         *pgxpool.Pool
	AuthService  *service.AuthService
	ProjectSvc   *service.ProjectService
	SuiteSvc     *service.TestSuiteService
	TestCaseSvc  *service.TestCaseService
	RunSvc       *service.TestRunService
	WorkerSvc    *service.WorkerService
	DashboardSvc *service.DashboardService
	OpenAPISvc   *service.OpenAPIService
	LoadRepo     service.LoadMetricsRepository
	Hub          *ws.Hub
}

func New(deps *Dependencies) *Server {
	r := chi.NewRouter()

	// Setup routes
	handlerDeps := &handler.Deps{
		AuthService:  deps.AuthService,
		ProjectSvc:   deps.ProjectSvc,
		SuiteSvc:     deps.SuiteSvc,
		TestCaseSvc:  deps.TestCaseSvc,
		RunSvc:       deps.RunSvc,
		WorkerSvc:    deps.WorkerSvc,
		DashboardSvc: deps.DashboardSvc,
		OpenAPISvc:   deps.OpenAPISvc,
		LoadRepo:     deps.LoadRepo,
		Hub:          deps.Hub,
	}

	SetupRoutes(r, handlerDeps, deps.Logger)

	httpServer := &http.Server{
		Addr:         deps.Config.Server.Addr(),
		Handler:      r,
		ReadTimeout:  deps.Config.Server.ReadTimeout,
		WriteTimeout: deps.Config.Server.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     deps.Logger,
		pool:       deps.Pool,
		hub:        deps.Hub,
	}
}

// Start begins listening and gracefully shuts down on SIGINT/SIGTERM.
func (s *Server) Start() error {
	// Start the WebSocket hub
	go s.hub.Run()

	// Channel to receive errors from the server
	errChan := make(chan error, 1)

	go func() {
		s.logger.Info("server starting", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
		s.logger.Info("server shutting down...")
	}

	// Graceful shutdown with 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to shutdown", "error", err)
		return err
	}

	if s.pool != nil {
		s.pool.Close()
	}

	s.logger.Info("server stopped gracefully")
	return nil
}
