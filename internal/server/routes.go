package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/qnguyenhong/automation-platform/internal/handler"
	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
)

func SetupRoutes(r *chi.Mux, deps *handler.Deps, logger *slog.Logger) {
	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Worker-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Post("/auth/login", deps.Login)
		r.Post("/auth/register", deps.Register)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.AuthService))

			// Auth
			r.Get("/auth/me", deps.GetCurrentUser)

			// Projects
			r.Route("/projects", func(r chi.Router) {
				r.Get("/", deps.ListProjects)
				r.Post("/", deps.CreateProject)

				r.Route("/{projectID}", func(r chi.Router) {
					r.Get("/", deps.GetProject)
					r.Put("/", deps.UpdateProject)
					r.Delete("/", deps.DeleteProject)

					// Environments under project
					r.Route("/environments", func(r chi.Router) {
						r.Get("/", deps.ListEnvironments)
						r.Post("/", deps.CreateEnvironment)
						r.Route("/{envID}", func(r chi.Router) {
							r.Get("/", deps.GetEnvironment)
							r.Put("/", deps.UpdateEnvironment)
							r.Delete("/", deps.DeleteEnvironment)
						})
					})

					// Datasets under project
					r.Route("/datasets", func(r chi.Router) {
						r.Get("/", deps.ListDatasets)
						r.Post("/", deps.CreateDataset)
					})

					// Webhooks under project
					r.Route("/webhooks", func(r chi.Router) {
						r.Get("/", deps.ListWebhooks)
						r.Post("/", deps.CreateWebhook)
					})

					// Suites under project
					r.Route("/suites", func(r chi.Router) {
						r.Get("/", deps.ListTestSuites)
						r.Post("/", deps.CreateTestSuite)

						r.Route("/{suiteID}", func(r chi.Router) {
							r.Get("/", deps.GetTestSuite)
							r.Put("/", deps.UpdateTestSuite)
							r.Delete("/", deps.DeleteTestSuite)

							// Cases under suite
							r.Route("/cases", func(r chi.Router) {
								r.Get("/", deps.ListTestCases)
								r.Post("/", deps.CreateTestCase)

								r.Route("/{caseID}", func(r chi.Router) {
									r.Get("/", deps.GetTestCase)
									r.Put("/", deps.UpdateTestCase)
									r.Delete("/", deps.DeleteTestCase)
								})
							})

							// Trigger run under suite
							r.Post("/runs", deps.TriggerRun)
						})
					})

					// Runs under project
					r.Get("/runs", deps.ListTestRuns)
				})
			})

			// Runs (global)
			r.Get("/runs", deps.ListAllTestRuns)
			r.Route("/runs/{runID}", func(r chi.Router) {
				r.Get("/", deps.GetTestRun)
				r.Post("/cancel", deps.CancelRun)
				r.Get("/results", deps.ListTestResults)
				r.Get("/load-metrics", deps.GetRunLoadMetrics)
				r.Route("/results/{resultID}", func(r chi.Router) {
					r.Get("/", deps.GetTestResult)
					r.Get("/load-metrics", deps.GetResultLoadMetrics)
				})
			})

			// Workers
			r.Route("/workers", func(r chi.Router) {
				r.Get("/", deps.ListWorkers)
				r.Get("/{workerID}", deps.GetWorker)
				r.Delete("/{workerID}", deps.DeleteWorker)
			})

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/config", deps.ListNotificationConfigs)
				r.Post("/config", deps.CreateNotificationConfig)
				r.Put("/config/{configID}", deps.UpdateNotificationConfig)
				r.Delete("/config/{configID}", deps.DeleteNotificationConfig)
			})

			// Dashboard
			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/summary", deps.GetDashboardSummary)
				r.Get("/trends", deps.GetDashboardTrends)
				r.Get("/flaky", deps.GetFlakyTests)
			})

			// OpenAPI
			r.Route("/openapi", func(r chi.Router) {
				r.Post("/parse", deps.ParseOpenAPI)
				r.Post("/import", deps.ImportOpenAPI)
			})

			// Environments (under project)
			// Note: Also accessible via /projects/{projectID}/environments in the project route

			// Datasets
			r.Route("/datasets/{datasetID}", func(r chi.Router) {
				r.Get("/", deps.GetDataset)
				r.Put("/", deps.UpdateDataset)
				r.Delete("/", deps.DeleteDataset)
			})

			// Webhooks (management under project routes, trigger is public)
			r.Route("/webhooks/{webhookID}", func(r chi.Router) {
				r.Get("/", deps.GetWebhook)
				r.Put("/", deps.UpdateWebhook)
				r.Delete("/", deps.DeleteWebhook)
			})

			// Export
			r.Post("/export", deps.ExportSuite)

			// WebSocket
			r.Get("/ws", deps.HandleWebSocket)
		})

		// Worker-facing endpoints (separate auth via worker token)
		r.Group(func(r chi.Router) {
			r.Use(middleware.WorkerAuth(deps.WorkerSvc))

			r.Post("/workers/register", deps.RegisterWorker)
			r.Put("/workers/{workerID}/heartbeat", deps.WorkerHeartbeat)
			r.Get("/jobs/next", deps.GetNextJob)
			r.Post("/jobs/{jobID}/result", deps.SubmitJobResult)
			r.Post("/jobs/{jobID}/log", deps.SubmitJobLog)
		})

		// Public webhook trigger (no auth required, uses secret)
		r.Post("/webhooks/trigger", deps.TriggerWebhook)
	})
}
