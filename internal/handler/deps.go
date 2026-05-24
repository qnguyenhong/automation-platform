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
	Hub           *ws.Hub
}
