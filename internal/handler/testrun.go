package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/server/middleware"
	"github.com/qnguyenhong/automation-platform/pkg/httputil"
)

func (d *Deps) TriggerRun(w http.ResponseWriter, r *http.Request) {
	suiteID, err := uuid.Parse(chi.URLParam(r, "suiteID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid suite ID")
		return
	}

	var req model.TriggerRunRequest
	if err := httputil.Decode(r, &req); err != nil {
		req.Trigger = model.TriggerManual
	}

	userID := middleware.GetUserID(r.Context())
	run, err := d.RunSvc.Trigger(r.Context(), suiteID, userID, req)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	// Notify via WebSocket
	d.Hub.Broadcast(map[string]any{
		"type": "run_started",
		"run":  run,
	})

	httputil.JSON(w, http.StatusCreated, run)
}

func (d *Deps) ListTestRuns(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid project ID")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	runs, total, err := d.RunSvc.ListByProject(r.Context(), projectID, page, perPage)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSONWithMeta(w, http.StatusOK, runs, httputil.Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: (int(total) + perPage - 1) / perPage,
	})
}

func (d *Deps) GetTestRun(w http.ResponseWriter, r *http.Request) {
	runID, err := uuid.Parse(chi.URLParam(r, "runID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	run, err := d.RunSvc.Get(r.Context(), runID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, run)
}

func (d *Deps) CancelRun(w http.ResponseWriter, r *http.Request) {
	runID, err := uuid.Parse(chi.URLParam(r, "runID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	if err := d.RunSvc.Cancel(r.Context(), runID); err != nil {
		httputil.AppError(w, err)
		return
	}

	d.Hub.Broadcast(map[string]any{
		"type":   "run_cancelled",
		"run_id": runID.String(),
	})

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (d *Deps) ListTestResults(w http.ResponseWriter, r *http.Request) {
	runID, err := uuid.Parse(chi.URLParam(r, "runID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid run ID")
		return
	}

	results, err := d.RunSvc.ListResults(r.Context(), runID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, results)
}

func (d *Deps) GetTestResult(w http.ResponseWriter, r *http.Request) {
	resultID, err := uuid.Parse(chi.URLParam(r, "resultID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid result ID")
		return
	}

	result, err := d.RunSvc.GetResult(r.Context(), resultID)
	if err != nil {
		httputil.AppError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}
