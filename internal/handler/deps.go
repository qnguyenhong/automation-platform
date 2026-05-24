package handler

import (
	"github.com/qnguyenhong/automation-platform/internal/service"
	"github.com/qnguyenhong/automation-platform/internal/ws"
)

// Deps holds all dependencies for HTTP handlers.
type Deps struct {
	AuthService   *service.AuthService
	ProjectSvc    *service.ProjectService
	SuiteSvc      *service.TestSuiteService
	TestCaseSvc   *service.TestCaseService
	RunSvc        *service.TestRunService
	WorkerSvc     *service.WorkerService
	DashboardSvc  *service.DashboardService
	OpenAPISvc    *service.OpenAPIService
	EnvSvc        *service.EnvironmentService
	DatasetSvc    *service.DatasetService
	WebhookSvc    *service.WebhookService
	ExportSvc     *service.ExportService
	LoadRepo      service.LoadMetricsRepository
	Hub           *ws.Hub
}
