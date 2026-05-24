package main

import (
	"context"
	"log"
	"os"

	"github.com/qnguyenhong/automation-platform/internal/repository"
	"github.com/qnguyenhong/automation-platform/internal/server"
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/internal/ws"
	"github.com/qnguyenhong/automation-platform/pkg/config"
	"github.com/qnguyenhong/automation-platform/pkg/database"
	"github.com/qnguyenhong/automation-platform/pkg/logger"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lg := logger.New(cfg.Log.Level, cfg.Log.Format)

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.Database.URL)
	if err != nil {
		lg.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	projectRepo := repository.NewProjectRepo(pool)
	suiteRepo := repository.NewTestSuiteRepo(pool)
	caseRepo := repository.NewTestCaseRepo(pool)
	runRepo := repository.NewTestRunRepo(pool)
	resultRepo := repository.NewTestResultRepo(pool)
	workerRepo := repository.NewWorkerRepo(pool)
	dashboardRepo := repository.NewDashboardRepo(pool)
	openapiRepo := repository.NewOpenAPIRepo(pool)
	loadRepo := repository.NewLoadMetricsRepo(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTExpiry)
	projectSvc := service.NewProjectService(projectRepo)
	suiteSvc := service.NewTestSuiteService(suiteRepo)
	caseSvc := service.NewTestCaseService(caseRepo)
	workerSvc := service.NewWorkerService(workerRepo, cfg.Auth.JWTSecret)
	dashboardSvc := service.NewDashboardService(dashboardRepo)

	dispatcher := service.NewDispatcher(lg, workerSvc, 100)
	workerSvc.SetDispatcher(dispatcher)
	workerSvc.StartHealthCheck(ctx, lg)
	runSvc := service.NewTestRunService(runRepo, resultRepo, suiteRepo, caseRepo, loadRepo, dispatcher)

	openapiSvc := service.NewOpenAPIService(openapiRepo, suiteSvc, caseSvc)

	hub := ws.NewHub(lg)

	srv := server.New(&server.Dependencies{
		Config:       cfg,
		Logger:       lg,
		Pool:         pool,
		AuthService:  authSvc,
		ProjectSvc:   projectSvc,
		SuiteSvc:     suiteSvc,
		TestCaseSvc:  caseSvc,
		RunSvc:       runSvc,
		WorkerSvc:    workerSvc,
		DashboardSvc: dashboardSvc,
		OpenAPISvc:   openapiSvc,
		LoadRepo:     loadRepo,
		Hub:          hub,
	})

	if err := srv.Start(); err != nil {
		lg.Error("server error", "error", err)
		os.Exit(1)
	}
}
